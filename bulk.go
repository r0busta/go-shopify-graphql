package shopify

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/r0busta/go-shopify-graphql/rand"
	"github.com/r0busta/go-shopify-graphql/utils"
	"github.com/shurcooL/graphql"
)

type BulkOperationService interface {
	BulkQuery(query string, v interface{}) error
}

type BulkOperationServiceOp struct {
	client *Client
}

type queryCurrentBulkOperation struct {
	CurrentBulkOperation currentBulkOperation
}

type currentBulkOperation struct {
	ID             graphql.ID
	Status         graphql.String
	ErrorCode      graphql.String
	CreatedAt      graphql.String
	CompletedAt    graphql.String
	ObjectCount    graphql.String
	FileSize       graphql.String
	URL            graphql.String
	PartialDataURL graphql.String
}

type bulkOperationRunQueryResult struct {
	BulkOperation struct {
		ID graphql.ID
	}
	UserErrors []struct {
		Field   []graphql.String
		Message graphql.String
	}
}

type mutationBulkOperationRunQuery struct {
	BulkOperationRunQueryResult bulkOperationRunQueryResult `graphql:"bulkOperationRunQuery(query: $query)"`
}

type bulkOperationCancelResult struct {
	BulkOperation struct {
		ID graphql.ID
	}
	UserErrors []struct {
		Field   []graphql.String
		Message graphql.String
	}
}

type mutationBulkOperationRunQueryCancel struct {
	BulkOperationCancelResult bulkOperationCancelResult `graphql:"bulkOperationCancel(id: $id)"`
}

func (s *BulkOperationServiceOp) postBulkQuery(query string) error {
	m := mutationBulkOperationRunQuery{}
	vars := map[string]interface{}{
		"query": graphql.String(query),
	}

	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}
	if len(m.BulkOperationRunQueryResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.BulkOperationRunQueryResult.UserErrors)
	}

	return nil
}

func (s *BulkOperationServiceOp) getBulkQueryResult() (url string, err error) {
	q := queryCurrentBulkOperation{}
	err = s.client.gql.Query(context.Background(), &q, nil)
	if err != nil {
		return
	}

	// Start polling the operation's status
	for q.CurrentBulkOperation.Status == "CREATED" || q.CurrentBulkOperation.Status == "RUNNING" {
		log.Println("Bulk operation still running...")
		time.Sleep(1 * time.Second)

		err = s.client.gql.Query(context.Background(), &q, nil)
		if err != nil {
			log.Printf("%+v", q)
			return
		}
	}
	log.Printf("Bulk operation finished with the status: %s", q.CurrentBulkOperation.Status)

	if q.CurrentBulkOperation.ErrorCode != "" {
		log.Printf("%+v", q)
		err = fmt.Errorf("Bulk operation error: %s", q.CurrentBulkOperation.ErrorCode)
		return
	}

	if q.CurrentBulkOperation.ObjectCount == "0" {
		return
	}

	url = string(q.CurrentBulkOperation.URL)
	return
}

func (s *BulkOperationServiceOp) cancelRunningBulkQuery() (err error) {
	q := queryCurrentBulkOperation{}

	err = s.client.gql.Query(context.Background(), &q, nil)
	if err != nil {
		return
	}

	if q.CurrentBulkOperation.Status == "RUNNING" {
		log.Println("Canceling running operation")
		operationID := q.CurrentBulkOperation.ID

		m := mutationBulkOperationRunQueryCancel{}
		vars := map[string]interface{}{
			"id": graphql.ID(operationID),
		}

		err = s.client.gql.Mutate(context.Background(), &m, vars)
		if err != nil {
			return err
		}
		if len(m.BulkOperationCancelResult.UserErrors) > 0 {
			return fmt.Errorf("%+v", m.BulkOperationCancelResult.UserErrors)
		}

		err = s.client.gql.Query(context.Background(), &q, nil)
		if err != nil {
			return
		}

		for q.CurrentBulkOperation.Status == "CANCELING" {
			log.Println("Bulk operation still canceling...")
			err = s.client.gql.Query(context.Background(), &q, nil)
			if err != nil {
				return
			}
		}
		log.Printf("Bulk operation cancelled")
	}

	return
}

func (s *BulkOperationServiceOp) BulkQuery(query string, out interface{}) (err error) {
	err = s.cancelRunningBulkQuery()
	if err != nil {
		return
	}

	err = s.postBulkQuery(query)
	if err != nil {
		return
	}

	url, err := s.getBulkQueryResult()
	if err != nil || url == "" {
		return
	}

	filename := fmt.Sprintf("%s%s", rand.String(10), ".jsonl")
	resultFile := filepath.Join(os.TempDir(), filename)
	err = utils.DownloadFile(resultFile, url)
	if err != nil {
		return
	}

	err = parseBulkQueryResult(resultFile, out)
	if err != nil {
		return
	}

	return
}

func parseBulkQueryResult(resultFile string, out interface{}) (err error) {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		err = fmt.Errorf("'records' is not a pointer")
		return
	}

	outValue := reflect.ValueOf(out)
	outSlice := outValue.Elem()
	if outSlice.Kind() != reflect.Slice {
		err = fmt.Errorf("'records' is not a  pointer to a slice interface")
		return
	}

	sliceItemType := outSlice.Type().Elem() // slice item type
	sliceItemKind := sliceItemType.Kind()
	itemType := sliceItemType // slice item underlying type
	if sliceItemKind == reflect.Ptr {
		itemType = itemType.Elem()
	}

	f, err := os.Open(resultFile)
	if err != nil {
		return
	}
	defer utils.CloseFile(f)

	reader := bufio.NewReader(f)
	json := jsoniter.ConfigFastest

	for {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}

		itemVal := reflect.New(itemType)
		err = json.Unmarshal(line, itemVal.Interface())
		if err != nil {
			return
		}

		if sliceItemKind == reflect.Ptr {
			outSlice.Set(reflect.Append(outSlice, itemVal))
		} else {
			outSlice.Set(reflect.Append(outSlice, itemVal.Elem()))
		}
	}

	if err != nil && err != io.EOF {
		return
	}

	err = nil
	return
}
