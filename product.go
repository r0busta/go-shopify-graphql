package shopify

import (
	"context"
	"fmt"
	"strings"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	"github.com/r0busta/graphql"
)

//go:generate mockgen -destination=./mock/product_service.go -package=mock . ProductService
type ProductService interface {
	List(query string) ([]model.Product, error)
	ListAll() ([]model.Product, error)

	Get(gid graphql.ID) (*model.Product, error)

	Create(product model.ProductInput, media []model.CreateMediaInput) (string, error)

	Update(product model.ProductInput) error

	Delete(product model.ProductDeleteInput) error
}

type ProductServiceOp struct {
	client *Client
}

type mutationProductCreate struct {
	ProductCreateResult model.PriceRuleCreatePayload `graphql:"productCreate(input: $input, media: $media)" json:"productCreate"`
}

type mutationProductUpdate struct {
	ProductUpdateResult model.ProductUpdatePayload `graphql:"productUpdate(input: $input)" json:"productUpdate"`
}

type mutationProductDelete struct {
	ProductDeleteResult model.ProductDeletePayload `graphql:"productDelete(input: $input)" json:"productDelete"`
}

const productBaseQuery = `
	id
	legacyResourceId
	handle
	options{
		name
		values
	}
	tags
	title
	description
	priceRangeV2{
		minVariantPrice{
			amount
			currencyCode
		}
		maxVariantPrice{
			amount
			currencyCode
		}
	}
	productType
	vendor
	totalInventory
	onlineStoreUrl	
	descriptionHtml
	seo{
		description
		title
	}
	templateSuffix
`

var productQuery = fmt.Sprintf(`
	%s
	variants(first:250, after: $cursor){
		edges{
			node{
				id
				legacyResourceId
				sku
				selectedOptions{
					name
					value
				}
				compareAtPrice
				price
				inventoryQuantity
				inventoryItem{
					id
					legacyResourceId							
				}
			}
		}
		pageInfo{
			hasNextPage
		}
	}
`, productBaseQuery)

var productBulkQuery = fmt.Sprintf(`
	%s
	metafields{
		edges{
			node{
				id
				legacyResourceId
				namespace
				key
				value
				valueType
			}
		}
	}
	variants{
		edges{
			node{
				id
				legacyResourceId
				sku
				selectedOptions{
					name
					value
				}
				compareAtPrice
				price
				inventoryQuantity
				inventoryItem{
					id
					legacyResourceId							
				}
			}
		}
	}
`, productBaseQuery)

func (s *ProductServiceOp) ListAll() ([]model.Product, error) {
	q := fmt.Sprintf(`
		{
			products{
				edges{
					node{
						%s
					}
				}
			}
		}
	`, productBulkQuery)

	res := []model.Product{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []model.Product{}, err
	}

	return res, nil
}

func (s *ProductServiceOp) List(query string) ([]model.Product, error) {
	q := fmt.Sprintf(`
		{
			products(query: "$query"){
				edges{
					node{
						%s
					}
				}
			}
		}
	`, productBulkQuery)

	q = strings.ReplaceAll(q, "$query", query)

	res := []model.Product{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []model.Product{}, err
	}

	return res, nil
}

func (s *ProductServiceOp) Get(id graphql.ID) (*model.Product, error) {
	out, err := s.getPage(id, "")
	if err != nil {
		return nil, err
	}

	nextPageData := out
	hasNextPage := out.Variants.PageInfo.HasNextPage
	for hasNextPage && len(nextPageData.Variants.Edges) > 0 {
		cursor := nextPageData.Variants.Edges[len(nextPageData.Variants.Edges)-1].Cursor
		nextPageData, err := s.getPage(id, cursor.String)
		if err != nil {
			return nil, err
		}
		out.Variants.Edges = append(out.Variants.Edges, nextPageData.Variants.Edges...)
		hasNextPage = nextPageData.Variants.PageInfo.HasNextPage
	}

	return out, nil
}

func (s *ProductServiceOp) getPage(id graphql.ID, cursor string) (*model.Product, error) {
	q := fmt.Sprintf(`
		query product($id: ID!, $cursor: String) {
			product(id: $id){
				%s
			}
		}
	`, productQuery)

	vars := map[string]interface{}{
		"id": id,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	out := struct {
		Product *model.Product `json:"product"`
	}{}
	err := s.client.gql.QueryString(context.Background(), q, vars, &out)
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) Create(product model.ProductInput, media []model.CreateMediaInput) (string, error) {
	m := mutationProductCreate{}

	vars := map[string]interface{}{
		"input": product,
		"media": media,
	}

	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return "", err
	}

	if len(m.ProductCreateResult.UserErrors) > 0 {
		return "", fmt.Errorf("%+v", m.ProductCreateResult.UserErrors)
	}

	return m.ProductCreateResult.PriceRule.ID.String, nil
}

func (s *ProductServiceOp) Update(product model.ProductInput) error {
	m := mutationProductUpdate{}

	vars := map[string]interface{}{
		"input": product,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.ProductUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductUpdateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) Delete(product model.ProductDeleteInput) error {
	m := mutationProductDelete{}

	vars := map[string]interface{}{
		"input": product,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.ProductDeleteResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductDeleteResult.UserErrors)
	}

	return nil
}
