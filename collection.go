package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	"github.com/r0busta/graphql"
	log "github.com/sirupsen/logrus"
)

type CollectionService interface {
	ListAll() ([]model.Collection, error)

	Get(id graphql.ID) (*model.Collection, error)

	Create(collection model.CollectionInput) (string, error)
	CreateBulk(collections []model.CollectionInput) error

	Update(collection model.CollectionInput) error
}

type CollectionServiceOp struct {
	client *Client
}

type mutationCollectionCreate struct {
	CollectionCreateResult model.CollectionCreatePayload `graphql:"collectionCreate(input: $input)" json:"collectionCreate"`
}

type mutationCollectionUpdate struct {
	CollectionCreateResult model.CollectionUpdatePayload `graphql:"collectionUpdate(input: $input)" json:"collectionUpdate"`
}

var collectionQuery = `
	id
	handle	
	title

	products(first:250, after: $cursor){
		edges{
			node{
				id
			}
			cursor
		}
		pageInfo{
			hasNextPage
		}		
	}	
`

var collectionBulkQuery = `
	id
	handle	
	title
`

func (s *CollectionServiceOp) ListAll() ([]model.Collection, error) {
	q := fmt.Sprintf(`
		{
			collections{
				edges{
					node{
						%s
					}
				}
			}
		}
	`, collectionBulkQuery)

	res := []model.Collection{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []model.Collection{}, err
	}

	return res, nil
}

func (s *CollectionServiceOp) Get(id graphql.ID) (*model.Collection, error) {
	out, err := s.getPage(id, "")
	if err != nil {
		return nil, err
	}

	nextPageData := out
	hasNextPage := out.Products.PageInfo.HasNextPage
	for hasNextPage && len(nextPageData.Products.Edges) > 0 {
		cursor := nextPageData.Products.Edges[len(nextPageData.Products.Edges)-1].Cursor
		nextPageData, err := s.getPage(id, cursor.String)
		if err != nil {
			return nil, err
		}
		out.Products.Edges = append(out.Products.Edges, nextPageData.Products.Edges...)
		hasNextPage = nextPageData.Products.PageInfo.HasNextPage
	}

	return out, nil
}

func (s *CollectionServiceOp) getPage(id graphql.ID, cursor string) (*model.Collection, error) {
	q := fmt.Sprintf(`
		query collection($id: ID!, $cursor: String) {
			collection(id: $id){
				%s
			}
		}
	`, collectionQuery)

	vars := map[string]interface{}{
		"id": id,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	out := struct {
		Collection *model.Collection `json:"collection"`
	}{}
	err := s.client.gql.QueryString(context.Background(), q, vars, &out)
	if err != nil {
		return nil, err
	}

	return out.Collection, nil
}

func (s *CollectionServiceOp) CreateBulk(collections []model.CollectionInput) error {
	for _, c := range collections {
		_, err := s.client.Collection.Create(c)
		if err != nil {
			log.Warnf("Couldn't create collection (%v): %s", c, err)
		}
	}

	return nil
}

func (s *CollectionServiceOp) Create(collection model.CollectionInput) (string, error) {
	m := mutationCollectionCreate{}

	vars := map[string]interface{}{
		"input": collection,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return "", err
	}

	if len(m.CollectionCreateResult.UserErrors) > 0 {
		return "", fmt.Errorf("%+v", m.CollectionCreateResult.UserErrors)
	}

	return m.CollectionCreateResult.Collection.ID.String, nil
}

func (s *CollectionServiceOp) Update(collection model.CollectionInput) error {
	m := mutationCollectionUpdate{}

	vars := map[string]interface{}{
		"input": collection,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.CollectionCreateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.CollectionCreateResult.UserErrors)
	}

	return nil
}
