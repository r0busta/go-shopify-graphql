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
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/r0busta/go-shopify-graphql/v2/rand"
	"github.com/r0busta/go-shopify-graphql/v2/utils"
	"github.com/r0busta/graphql"
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
	ID             graphql.ID     `json:"id"`
	Status         graphql.String `json:"status"`
	ErrorCode      graphql.String `json:"errorCode"`
	CreatedAt      graphql.String `json:"createdAt"`
	CompletedAt    graphql.String `json:"completedAt"`
	ObjectCount    graphql.String `json:"objectCount"`
	FileSize       graphql.String `json:"fileSize"`
	URL            graphql.String `json:"url"`
	PartialDataURL graphql.String `json:"partialDataUrl"`
	Query          graphql.String `json:"query"`
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

var gidRegex *regexp.Regexp

func init() {
	gidRegex = regexp.MustCompile(`^gid://shopify/(\w+)/\d+$`)
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

	if q.CurrentBulkOperation.Status != "COMPLETED" {
		log.Printf("%+v", q)
		err = fmt.Errorf("Bulk operation didn't complete, status=%s", q.CurrentBulkOperation.Status)
		return
	}

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

	childrenLookup := make(map[string]interface{})

	for {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}

		parentID := json.Get(line, "__parentId")
		if parentID.LastError() == nil {
			gid := json.Get(line, "id")
			if gid.LastError() != nil {
				return fmt.Errorf("Connection type must query `id` field")
			}
			childObjType, childrenFieldName, err := concludeObjectType(gid.ToString())
			if err != nil {
				return err
			}
			childItem := reflect.New(childObjType).Interface()
			err = json.Unmarshal(line, &childItem)
			if err != nil {
				return err
			}
			childItemVal := reflect.ValueOf(childItem).Elem()

			var childrenSlice reflect.Value
			var children map[string]interface{}
			if val, ok := childrenLookup[parentID.ToString()]; ok {
				children = val.(map[string]interface{})
			} else {
				children = make(map[string]interface{})
			}

			if val, ok := children[childrenFieldName]; ok {
				childrenSlice = reflect.ValueOf(val)
			} else {
				childrenSlice = reflect.MakeSlice(reflect.SliceOf(childObjType), 0, 10)
			}

			childrenSlice = reflect.Append(childrenSlice, childItemVal)

			children[childrenFieldName] = childrenSlice.Interface()
			childrenLookup[parentID.ToString()] = children

			continue
		}

		item := reflect.New(itemType).Interface()
		err = json.Unmarshal(line, &item)
		if err != nil {
			return
		}
		itemVal := reflect.ValueOf(item)

		if sliceItemKind == reflect.Ptr {
			outSlice.Set(reflect.Append(outSlice, itemVal))
		} else {
			outSlice.Set(reflect.Append(outSlice, itemVal.Elem()))
		}
	}

	if len(childrenLookup) > 0 {
		for i := 0; i < outSlice.Len(); i++ {
			parent := outSlice.Index(i).Elem()
			parentIDField := parent.FieldByName("ID")
			if parentIDField.IsZero() {
				return fmt.Errorf("No ID field on the first level")
			}
			parentID := parentIDField.Interface().(string)
			if children, ok := childrenLookup[parentID]; ok {
				childrenVal := reflect.ValueOf(children)
				iter := childrenVal.MapRange()
				for iter.Next() {
					k := iter.Key()
					v := reflect.ValueOf(iter.Value().Interface())
					parent.FieldByName(k.String()).Set(v)
				}
			}
		}
	}

	if err != nil && err != io.EOF {
		return
	}

	err = nil
	return
}

func concludeObjectType(gid string) (reflect.Type, string, error) {
	submatches := gidRegex.FindStringSubmatch(gid)
	if len(submatches) != 2 {
		return reflect.TypeOf(nil), "", fmt.Errorf("malformed gid=`%s`", gid)
	}
	resource := submatches[1]
	switch resource {
	case "LineItem":
		return reflect.TypeOf(LineItem{}), fmt.Sprintf("%ss", resource), nil
	case "Metafield":
		return reflect.TypeOf(Metafield{}), fmt.Sprintf("%ss", resource), nil
	case "Order":
		return reflect.TypeOf(Order{}), fmt.Sprintf("%ss", resource), nil
	case "Product":
		return reflect.TypeOf(ProductBulkResult{}), fmt.Sprintf("%ss", resource), nil
	case "ProductVariant":
		return reflect.TypeOf(ProductVariant{}), fmt.Sprintf("%ss", resource), nil
	default:
		return reflect.TypeOf(nil), "", fmt.Errorf("`%s` not implemented type", resource)
	}
}
