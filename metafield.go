package shopify

import (
	"context"
	"fmt"
	"strings"

	"github.com/r0busta/go-shopify-graphql-model/v5/graph/model"
)

//go:generate mockgen -destination=./mock/metafield_service.go -package=mock . MetafieldService
type MetafieldService interface {
	ListAllShopMetafields(ctx context.Context) ([]model.Metafield, error)
	ListShopMetafieldsByNamespace(ctx context.Context, namespace string) ([]model.Metafield, error)

	GetShopMetafieldByKey(ctx context.Context, namespace, key string) (*model.Metafield, error)

	Delete(ctx context.Context, metafield model.MetafieldIdentifierInput) error
	DeleteBulk(ctx context.Context, metafield []model.MetafieldIdentifierInput) error
}

type MetafieldServiceOp struct {
	client *Client
}

var _ MetafieldService = &MetafieldServiceOp{}

type mutationMetafieldsDelete struct {
	MetafieldsDeleteResult struct {
		DeletedMetafields []*model.MetafieldIdentifier `json:"deletedMetafields,omitempty"`
		UserErrors        []model.UserError            `json:"userErrors,omitempty"`
	} `graphql:"metafieldsDelete(metafields: $metafields)" json:"metafieldsDelete"`
}

func (s *MetafieldServiceOp) ListAllShopMetafields(ctx context.Context) ([]model.Metafield, error) {
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
							type
						}
					}
				}	  
			}
		}
`

	res := []model.Metafield{}
	err := s.client.BulkOperation.BulkQuery(ctx, q, &res)
	if err != nil {
		return nil, fmt.Errorf("bulk query: %w", err)
	}

	return res, nil
}

func (s *MetafieldServiceOp) ListShopMetafieldsByNamespace(ctx context.Context, namespace string) ([]model.Metafield, error) {
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
							type
						}
					}
				}	  
			}
		}
`
	q = strings.ReplaceAll(q, "$namespace", namespace)

	res := []model.Metafield{}
	err := s.client.BulkOperation.BulkQuery(ctx, q, &res)
	if err != nil {
		return nil, fmt.Errorf("bulk query: %w", err)
	}

	return res, nil
}

// GetShopMetafieldByKey returns the shop metafield with the given namespace
// and key, or nil if there is none.
func (s *MetafieldServiceOp) GetShopMetafieldByKey(ctx context.Context, namespace, key string) (*model.Metafield, error) {
	q := `
		query shopMetafield($namespace: String, $key: String!) {
			shop {
				metafield(namespace: $namespace, key: $key) {
					createdAt
					description
					id
					key
					legacyResourceId
					namespace
					ownerType
					updatedAt
					value
					type
				}
			}
		}
`
	vars := map[string]interface{}{
		"namespace": namespace,
		"key":       key,
	}

	var out struct {
		Shop struct {
			Metafield *model.Metafield `json:"metafield"`
		} `json:"shop"`
	}
	err := s.client.gql.QueryString(ctx, q, vars, &out)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out.Shop.Metafield, nil
}

// DeleteBulk deletes the given metafields in one metafieldsDelete mutation.
// Identifiers that match no metafield are skipped without an error.
func (s *MetafieldServiceOp) DeleteBulk(ctx context.Context, metafields []model.MetafieldIdentifierInput) error {
	_, err := s.deleteMetafields(ctx, metafields)
	return err
}

// Delete deletes one metafield and returns an error if it does not exist.
func (s *MetafieldServiceOp) Delete(ctx context.Context, metafield model.MetafieldIdentifierInput) error {
	deleted, err := s.deleteMetafields(ctx, []model.MetafieldIdentifierInput{metafield})
	if err != nil {
		return err
	}
	if len(deleted) == 0 || deleted[0] == nil {
		return fmt.Errorf("metafield %s.%s on %s not found", metafield.Namespace, metafield.Key, metafield.OwnerID)
	}

	return nil
}

// deleteMetafields returns one entry per identifier, nil where no metafield
// matched.
func (s *MetafieldServiceOp) deleteMetafields(ctx context.Context, metafields []model.MetafieldIdentifierInput) ([]*model.MetafieldIdentifier, error) {
	if len(metafields) == 0 {
		return nil, nil
	}

	m := mutationMetafieldsDelete{}

	vars := map[string]interface{}{
		"metafields": metafields,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.MetafieldsDeleteResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.MetafieldsDeleteResult.UserErrors)
	}

	return m.MetafieldsDeleteResult.DeletedMetafields, nil
}
