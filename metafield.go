package shopify

import (
	"context"
	"fmt"
	"strings"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	log "github.com/sirupsen/logrus"
)

type MetafieldService interface {
	ListAllShopMetafields() ([]*model.Metafield, error)
	ListShopMetafieldsByNamespace(namespace string) ([]*model.Metafield, error)

	GetShopMetafieldByKey(namespace, key string) (model.Metafield, error)

	Delete(metafield *model.MetafieldDeleteInput) error
	DeleteBulk(metafield []*model.MetafieldDeleteInput) error
}

type MetafieldServiceOp struct {
	client *Client
}
type mutationMetafieldDelete struct {
	MetafieldDeleteResult model.MetafieldDeletePayload `graphql:"metafieldDelete(input: $input)" json:"metafieldDelete"`
}

func (s *MetafieldServiceOp) ListAllShopMetafields() ([]*model.Metafield, error) {
	q := `
		{
			shop{
				metafields{
					edges{
						node{
							createdAt
							description
							id
							key
							legacyResourceId
							namespace
							ownerType
							updatedAt
							value
							valueType
						}
					}
				}	  
			}
		}
`

	res := []*model.Metafield{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*model.Metafield{}, err
	}

	return res, nil
}

func (s *MetafieldServiceOp) ListShopMetafieldsByNamespace(namespace string) ([]*model.Metafield, error) {
	q := `
		{
			shop{
				metafields(namespace: "$namespace"){
					edges{
						node{
							createdAt
							description
							id
							key
							legacyResourceId
							namespace
							ownerType
							updatedAt
							value
							valueType
						}
					}
				}	  
			}
		}
`
	q = strings.ReplaceAll(q, "$namespace", namespace)

	res := []*model.Metafield{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*model.Metafield{}, err
	}

	return res, nil
}

func (s *MetafieldServiceOp) GetShopMetafieldByKey(namespace, key string) (model.Metafield, error) {
	var q struct {
		Shop struct {
			Metafield model.Metafield `graphql:"metafield(namespace: $namespace, key: $key)"`
		} `graphql:"shop"`
	}
	vars := map[string]interface{}{
		"namespace": namespace,
		"key":       key,
	}

	err := s.client.gql.Query(context.Background(), &q, vars)
	if err != nil {
		return model.Metafield{}, err
	}

	return q.Shop.Metafield, nil
}

func (s *MetafieldServiceOp) DeleteBulk(metafields []*model.MetafieldDeleteInput) error {
	for _, m := range metafields {
		err := s.Delete(m)
		if err != nil {
			log.Warnf("Couldn't delete metafield (%v): %s", m, err)
		}
	}

	return nil
}

func (s *MetafieldServiceOp) Delete(metafield *model.MetafieldDeleteInput) error {
	m := mutationMetafieldDelete{}

	vars := map[string]interface{}{
		"input": metafield,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.MetafieldDeleteResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.MetafieldDeleteResult.UserErrors)
	}

	return nil
}
