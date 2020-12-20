package shopify

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/r0busta/graphql"
)

type MetafieldService interface {
	ListAllShopMetafields() ([]*Metafield, error)
	ListShopMetafieldsByNamespace(namespace string) ([]*Metafield, error)

	GetShopMetafieldByKey(namespace, key string) (Metafield, error)

	Delete(metafield MetafieldDeleteInput) error
	DeleteBulk(metafield []MetafieldDeleteInput) error
}

type MetafieldServiceOp struct {
	client *Client
}

type Metafield struct {
	// The date and time when the metafield was created.
	CreatedAt DateTime `json:"createdAt,omitempty"`
	// The description of a metafield.
	Description graphql.String `json:"description,omitempty"`
	// Globally unique identifier.
	ID graphql.ID `json:"id,omitempty"`
	// The key name for a metafield.
	Key graphql.String `json:"key,omitempty"`
	// The ID of the corresponding resource in the REST Admin API.
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
	// The namespace for a metafield.
	Namespace graphql.String `json:"namespace,omitempty"`
	// Owner type of a metafield visible to the Storefront API.
	OwnerType graphql.String `json:"ownerType,omitempty"`
	// The date and time when the metafield was updated.
	UpdatedAt DateTime `json:"updatedAt,omitempty"`
	// The value of a metafield.
	Value graphql.String `json:"value,omitempty"`
	// Represents the metafield value type.
	ValueType MetafieldValueType `json:"valueType,omitempty"`
}

type MetafieldDeleteInput struct {
	// The ID of the metafield to delete.
	ID graphql.ID `json:"id,omitempty"`
}

type mutationMetafieldDelete struct {
	MetafieldDeleteResult metafieldDeleteResult `graphql:"metafieldDelete(input: $input)" json:"metafieldDelete"`
}

type metafieldDeleteResult struct {
	DeletedID  string       `json:"deletedId,omitempty"`
	UserErrors []UserErrors `json:"userErrors"`
}

func (s *MetafieldServiceOp) ListAllShopMetafields() ([]*Metafield, error) {
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

	res := []*Metafield{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*Metafield{}, err
	}

	return res, nil
}

func (s *MetafieldServiceOp) ListShopMetafieldsByNamespace(namespace string) ([]*Metafield, error) {
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

	res := []*Metafield{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*Metafield{}, err
	}

	return res, nil
}

func (s *MetafieldServiceOp) GetShopMetafieldByKey(namespace, key string) (Metafield, error) {
	var q struct {
		Shop struct {
			Metafield Metafield `graphql:"metafield(namespace: $namespace, key: $key)"`
		} `graphql:"shop"`
	}
	vars := map[string]interface{}{
		"namespace": graphql.String(namespace),
		"key":       graphql.String(key),
	}

	err := s.client.gql.Query(context.Background(), &q, vars)
	if err != nil {
		return Metafield{}, err
	}

	return q.Shop.Metafield, nil
}

func (s *MetafieldServiceOp) DeleteBulk(metafields []MetafieldDeleteInput) error {
	for _, m := range metafields {
		err := s.Delete(m)
		if err != nil {
			log.Printf("Warning! Couldn't delete metafield (%v): %s", m, err)
		}
	}

	return nil
}

func (s *MetafieldServiceOp) Delete(metafield MetafieldDeleteInput) error {
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
