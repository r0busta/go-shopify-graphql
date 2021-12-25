package shopify

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/r0busta/go-shopify-graphql-model/v2/graph/model"
	"github.com/r0busta/go-shopify-graphql/v5/rand"
	"github.com/r0busta/go-shopify-graphql/v5/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

const (
	edgesFieldName = "Edges"
	nodeFieldName  = "Node"
)

//go:generate mockgen -destination=./mock/bulk_service.go -package=mock . BulkOperationService
type BulkOperationService interface {
	BulkQuery(query string, v interface{}) error

	PostBulkQuery(query string) (*string, error)
	GetCurrentBulkQuery() (*model.BulkOperation, error)
	GetCurrentBulkQueryResultURL() (*string, error)
	WaitForCurrentBulkQuery(interval time.Duration) (*model.BulkOperation, error)
	ShouldGetBulkQueryResultURL(id *string) (*string, error)
	CancelRunningBulkQuery() error
}

type BulkOperationServiceOp struct {
	client *Client
}

var _ BulkOperationService = &BulkOperationServiceOp{}

type mutationBulkOperationRunQuery struct {
	BulkOperationRunQueryResult model.BulkOperationRunQueryPayload `graphql:"bulkOperationRunQuery(query: $query)" json:"bulkOperationRunQuery"`
}

type mutationBulkOperationRunQueryCancel struct {
	BulkOperationCancelResult model.BulkOperationCancelPayload `graphql:"bulkOperationCancel(id: $id)" json:"bulkOperationCancel"`
}

var gidRegex *regexp.Regexp

func init() {
	gidRegex = regexp.MustCompile(`^gid://shopify/(\w+)/\d+$`)
}

func (s *BulkOperationServiceOp) PostBulkQuery(query string) (*string, error) {
	m := mutationBulkOperationRunQuery{}
	vars := map[string]interface{}{
		"query": null.StringFrom(query),
	}

	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return nil, fmt.Errorf("error posting bulk query: %s", err)
	}
	if len(m.BulkOperationRunQueryResult.UserErrors) > 0 {
		errors, _ := json.MarshalIndent(m.BulkOperationRunQueryResult.UserErrors, "", "    ")
		return nil, fmt.Errorf("error posting bulk query: %s", errors)
	}

	return &m.BulkOperationRunQueryResult.BulkOperation.ID, nil
}

func (s *BulkOperationServiceOp) GetCurrentBulkQuery() (*model.BulkOperation, error) {
	var q struct {
		CurrentBulkOperation struct {
			model.BulkOperation
		}
	}
	err := s.client.gql.Query(context.Background(), &q, nil)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return &q.CurrentBulkOperation.BulkOperation, nil
}

func (s *BulkOperationServiceOp) GetCurrentBulkQueryResultURL() (*string, error) {
	return s.ShouldGetBulkQueryResultURL(nil)
}

func (s *BulkOperationServiceOp) ShouldGetBulkQueryResultURL(id *string) (*string, error) {
	q, err := s.GetCurrentBulkQuery()
	if err != nil {
		return nil, fmt.Errorf("error getting current bulk operation: %s", err)
	}

	if id != nil && q.ID != *id {
		return nil, fmt.Errorf("Bulk operation ID doesn't match, got=%v, want=%v", q.ID, id)
	}

	q, _ = s.WaitForCurrentBulkQuery(1 * time.Second)
	if q.Status != "COMPLETED" {
		return nil, fmt.Errorf("Bulk operation didn't complete, status=%s, error_code=%s", q.Status, q.ErrorCode)
	}

	if q.ErrorCode != nil && q.ErrorCode.String() != "" {
		return nil, fmt.Errorf("Bulk operation error: %s", q.ErrorCode)
	}

	if q.ObjectCount == "0" {
		return nil, nil
	}

	if q.URL == nil {
		return nil, fmt.Errorf("empty URL result")
	}

	return q.URL, nil
}

func (s *BulkOperationServiceOp) WaitForCurrentBulkQuery(interval time.Duration) (*model.BulkOperation, error) {
	q, err := s.GetCurrentBulkQuery()
	if err != nil {
		return q, fmt.Errorf("CurrentBulkOperation query error: %s", err)
	}

	for q.Status == "CREATED" || q.Status == "RUNNING" || q.Status == "CANCELING" {
		log.Debugf("Bulk operation is still %s...", q.Status)
		time.Sleep(interval)

		q, err = s.GetCurrentBulkQuery()
		if err != nil {
			return q, fmt.Errorf("CurrentBulkOperation query error: %s", err)
		}
	}
	log.Debugf("Bulk operation ready, latest status=%s", q.Status)

	return q, nil
}

func (s *BulkOperationServiceOp) CancelRunningBulkQuery() error {
	q, err := s.GetCurrentBulkQuery()
	if err != nil {
		return err
	}

	if q.Status == "CREATED" || q.Status == "RUNNING" {
		log.Debugln("Canceling running operation")
		operationID := q.ID

		m := mutationBulkOperationRunQueryCancel{}
		vars := map[string]interface{}{
			"id": operationID,
		}

		err = s.client.gql.Mutate(context.Background(), &m, vars)
		if err != nil {
			return fmt.Errorf("mutation: %w", err)
		}
		if len(m.BulkOperationCancelResult.UserErrors) > 0 {
			return fmt.Errorf("%+v", m.BulkOperationCancelResult.UserErrors)
		}

		q, err = s.GetCurrentBulkQuery()
		if err != nil {
			return err
		}
		for q.Status == "CREATED" || q.Status == "RUNNING" || q.Status == "CANCELING" {
			log.Tracef("Bulk operation still %s...", q.Status)
			q, err = s.GetCurrentBulkQuery()
			if err != nil {
				return fmt.Errorf("get current bulk query: %w", err)
			}
		}
		log.Debugln("Bulk operation cancelled")
	}

	return nil
}

func (s *BulkOperationServiceOp) BulkQuery(query string, out interface{}) error {
	_, err := s.WaitForCurrentBulkQuery(1 * time.Second)
	if err != nil {
		return err
	}

	id, err := s.PostBulkQuery(query)
	if err != nil {
		return fmt.Errorf("post bulk query: %w", err)
	}

	if id == nil {
		return fmt.Errorf("Posted operation ID is nil")
	}

	url, err := s.ShouldGetBulkQueryResultURL(id)
	if err != nil {
		return fmt.Errorf("get bulk query result URL: %w", err)
	}

	if url == nil || *url == "" {
		return fmt.Errorf("Operation result URL is empty")
	}

	filename := fmt.Sprintf("%s%s", rand.String(10), ".jsonl")
	resultFile := filepath.Join(os.TempDir(), filename)
	err = utils.DownloadFile(resultFile, *url)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	err = parseBulkQueryResult(resultFile, out)
	if err != nil {
		return fmt.Errorf("parse bulk query result: %w", err)
	}

	return nil
}

func parseBulkQueryResult(resultFilePath string, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("the out arg is not a pointer")
	}

	outValue := reflect.ValueOf(out)
	outSlice := outValue.Elem()
	if outSlice.Kind() != reflect.Slice {
		return fmt.Errorf("the out arg is not a pointer to a slice interface")
	}

	sliceItemType := outSlice.Type().Elem() // slice item type
	sliceItemKind := sliceItemType.Kind()
	itemType := sliceItemType // slice item underlying type
	if sliceItemKind == reflect.Ptr {
		itemType = itemType.Elem()
	}

	resultPath, err := os.Open(resultFilePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer utils.CloseFile(resultPath)

	reader := bufio.NewReader(resultPath)
	json := jsoniter.ConfigFastest

	connectionSink := make(map[string]interface{})

	for {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}

		parentIDNode := json.Get(line, "__parentId")
		if parentIDNode.LastError() == nil {
			parentID := parentIDNode.ToString()

			gid := json.Get(line, "id")
			if gid.LastError() != nil {
				return fmt.Errorf("The connection type must query the `id` field")
			}
			edgeType, nodeType, connectionFieldName, err := concludeObjectType(gid.ToString())
			if err != nil {
				return err
			}
			node := reflect.New(nodeType).Interface()
			err = json.Unmarshal(line, &node)
			if err != nil {
				return fmt.Errorf("unmarshalling: %w", err)
			}
			nodeVal := reflect.ValueOf(node).Elem()

			var edge interface{}
			var edgeVal reflect.Value
			var nodeField reflect.Value
			if edgeType.Kind() == reflect.Ptr {
				edge = reflect.New(edgeType.Elem()).Interface()
				nodeField = reflect.ValueOf(edge).Elem().FieldByName(nodeFieldName)
				edgeVal = reflect.ValueOf(edge)
			} else {
				edge = reflect.New(edgeType).Interface()

				if reflect.ValueOf(edge).Kind() == reflect.Ptr {
					nodeField = reflect.ValueOf(edge).Elem().FieldByName(nodeFieldName)
				} else {
					nodeField = reflect.ValueOf(edge).FieldByName(nodeFieldName)
				}

				edgeVal = reflect.ValueOf(edge).Elem()
			}

			if !nodeField.IsValid() {
				return fmt.Errorf("Edge in the '%s' doesn't have the Node field", connectionFieldName)
			}
			nodeField.Set(nodeVal)

			var edgesSlice reflect.Value
			var edges map[string]interface{}
			if val, ok := connectionSink[parentID]; ok {
				var ok2 bool
				if edges, ok2 = val.(map[string]interface{}); !ok2 {
					return fmt.Errorf("The connection sink for parent ID '%s' is not a map", parentID)
				}
			} else {
				edges = make(map[string]interface{})
			}

			if val, ok := edges[connectionFieldName]; ok {
				edgesSlice = reflect.ValueOf(val)
			} else {
				edgesSlice = reflect.MakeSlice(reflect.SliceOf(edgeType), 0, 50)
			}

			edgesSlice = reflect.Append(edgesSlice, edgeVal)

			edges[connectionFieldName] = edgesSlice.Interface()
			connectionSink[parentID] = edges

			continue
		}

		item := reflect.New(itemType).Interface()
		err = json.Unmarshal(line, &item)
		if err != nil {
			return fmt.Errorf("unmarshalling: %w", err)
		}
		itemVal := reflect.ValueOf(item)

		if sliceItemKind == reflect.Ptr {
			outSlice.Set(reflect.Append(outSlice, itemVal))
		} else {
			outSlice.Set(reflect.Append(outSlice, itemVal.Elem()))
		}
	}

	if len(connectionSink) > 0 {
		for i := 0; i < outSlice.Len(); i++ {
			parent := outSlice.Index(i)
			if parent.Kind() == reflect.Ptr {
				parent = parent.Elem()
			}

			parentIDField := parent.FieldByName("ID")
			if parentIDField.IsZero() {
				return fmt.Errorf("No ID field on the first level")
			}

			var parentID string
			var ok bool
			if parentID, ok = parentIDField.Interface().(string); !ok {
				return fmt.Errorf("ID field on the first level is not a string")
			}

			if connection, ok := connectionSink[parentID]; ok {
				edgeVal := reflect.ValueOf(connection)
				iter := edgeVal.MapRange()
				for iter.Next() {
					connectionName := iter.Key()
					connectionField := parent.FieldByName(connectionName.String())
					if !connectionField.IsValid() {
						return fmt.Errorf("Connection '%s' is not defined on the parent type %s", connectionName.String(), parent.Type().String())
					}

					var connectionValue reflect.Value
					var edgesField reflect.Value
					if connectionField.Kind() == reflect.Ptr {
						connectionValue = reflect.ValueOf(reflect.New(connectionField.Type().Elem()).Interface())
						edgesField = connectionValue.Elem().FieldByName(edgesFieldName)
					} else {
						connectionValue = reflect.ValueOf(reflect.New(connectionField.Type()).Interface())
						edgesField = connectionValue.Elem().FieldByName(edgesFieldName)
					}

					if !edgesField.IsValid() {
						return fmt.Errorf("Connection %s in the '%s' doesn't have the Edges field", connectionName.String(), parent.Type().String())
					}

					edges := reflect.ValueOf(iter.Value().Interface())
					edgesField.Set(edges)

					connectionField.Set(connectionValue)
				}
			}
		}
	}

	// check if ReadBytes returned an error different from EOF
	if err != nil && err != io.EOF {
		return fmt.Errorf("reading the result file: %w", err)
	}

	return nil
}

func concludeObjectType(gid string) (reflect.Type, reflect.Type, string, error) {
	submatches := gidRegex.FindStringSubmatch(gid)
	if len(submatches) != 2 {
		return reflect.TypeOf(nil), reflect.TypeOf(nil), "", fmt.Errorf("malformed gid=`%s`", gid)
	}
	resource := submatches[1]
	switch resource {
	case "LineItem":
		return reflect.TypeOf(model.LineItemEdge{}), reflect.TypeOf(&model.LineItem{}), fmt.Sprintf("%ss", resource), nil
	case "FulfillmentOrderLineItem":
		return reflect.TypeOf(model.FulfillmentOrderLineItemEdge{}), reflect.TypeOf(&model.FulfillmentOrderLineItem{}), "LineItems", nil
	case "Metafield":
		return reflect.TypeOf(model.MetafieldEdge{}), reflect.TypeOf(&model.Metafield{}), fmt.Sprintf("%ss", resource), nil
	case "Order":
		return reflect.TypeOf(model.OrderEdge{}), reflect.TypeOf(&model.Order{}), fmt.Sprintf("%ss", resource), nil
	case "Product":
		return reflect.TypeOf(model.ProductEdge{}), reflect.TypeOf(&model.Product{}), fmt.Sprintf("%ss", resource), nil
	case "ProductVariant":
		return reflect.TypeOf(model.ProductVariantEdge{}), reflect.TypeOf(&model.ProductVariant{}), "Variants", nil
	default:
		return reflect.TypeOf(nil), reflect.TypeOf(nil), "", fmt.Errorf("`%s` not implemented type", resource)
	}
}
