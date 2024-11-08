package shopify

import (
	"context"
	"fmt"
	"strings"

	"github.com/r0busta/go-shopify-graphql-model/v4/graph/model"
)

//go:generate mockgen -destination=./mock/product_service.go -package=mock . ProductService
type ProductService interface {
	List(ctx context.Context, query string) ([]model.Product, error)
	ListAll(ctx context.Context) ([]model.Product, error)

	Get(ctx context.Context, id string) (*model.Product, error)

	Create(ctx context.Context, product model.ProductCreateInput, media []model.CreateMediaInput) (*string, error)

	Update(ctx context.Context, product model.ProductUpdateInput, media []model.CreateMediaInput) error

	Delete(ctx context.Context, product model.ProductDeleteInput) error

	VariantsBulkCreate(ctx context.Context, id string, input []model.ProductVariantsBulkInput, strategy model.ProductVariantsBulkCreateStrategy) error
	VariantsBulkUpdate(ctx context.Context, id string, input []model.ProductVariantsBulkInput) error
	VariantsBulkReorder(ctx context.Context, id string, input []model.ProductVariantPositionInput) error
}

type ProductServiceOp struct {
	client *Client
}

var _ ProductService = &ProductServiceOp{}

type mutationProductCreate struct {
	ProductCreateResult struct {
		Product *struct {
			ID string `json:"id,omitempty"`
		} `json:"product,omitempty"`

		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productCreate(product: $product, media: $media)" json:"productCreate"`
}

type mutationProductUpdate struct {
	ProductUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productUpdate(product: $product, media: $media)" json:"productUpdate"`
}

type mutationProductDelete struct {
	ProductDeleteResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productDelete(input: $input)" json:"productDelete"`
}

type mutationProductVariantsBulkCreate struct {
	ProductVariantsBulkCreateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productVariantsBulkCreate(productId: $productId, variants: $variants, strategy: $strategy)" json:"productVariantsBulkCreate"`
}

type mutationProductVariantsBulkUpdate struct {
	ProductVariantsBulkUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productVariantsBulkUpdate(productId: $productId, variants: $variants)" json:"productVariantsBulkUpdate"`
}

type mutationProductVariantsBulkReorder struct {
	ProductVariantsBulkReorderResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"productVariantsBulkReorder(positions: $positions, productId: $productId)" json:"productVariantsBulkReorder"`
}

const productBaseQuery = `
	id
	legacyResourceId
	handle
	options{
		name
		values
		position
	}
	tags
	title
	description
	descriptionPlainSummary
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
	customProductType
	featuredImage{
		id
		altText
		height
		width
		url
	}
`

var productQuery = fmt.Sprintf(`
	%s
	variants(first:100, after: $cursor){
		edges{
			node{
				id
				legacyResourceId
				title
				displayName
				sku
				selectedOptions{
					name
					value
				}
				position
				image {
					id
					altText
					height
					width
					url
				}
				compareAtPrice
				price
				inventoryQuantity
				inventoryItem{
					id
					legacyResourceId							
				}
				availableForSale
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
				type
			}
		}
	}
	variants{
		edges{
			node{
				id
				legacyResourceId
				title
				displayName
				sku
				selectedOptions{
					name
					value
				}
				position
				image {
					id
					altText
					height
					width
					url
				}
				compareAtPrice
				price
				inventoryQuantity
				inventoryItem{
					id
					legacyResourceId							
				}
				availableForSale
			}
		}
	}
`, productBaseQuery)

func (s *ProductServiceOp) ListAll(ctx context.Context) ([]model.Product, error) {
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
	err := s.client.BulkOperation.BulkQuery(ctx, q, &res)
	if err != nil {
		return []model.Product{}, err
	}

	return res, nil
}

func (s *ProductServiceOp) List(ctx context.Context, query string) ([]model.Product, error) {
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
	err := s.client.BulkOperation.BulkQuery(ctx, q, &res)
	if err != nil {
		return nil, fmt.Errorf("bulk query: %w", err)
	}

	return res, nil
}

func (s *ProductServiceOp) Get(ctx context.Context, id string) (*model.Product, error) {
	out, err := s.getPage(ctx, id, "")
	if err != nil {
		return nil, err
	}

	nextPageData := out
	hasNextPage := out.Variants.PageInfo.HasNextPage
	for hasNextPage && len(nextPageData.Variants.Edges) > 0 {
		cursor := nextPageData.Variants.Edges[len(nextPageData.Variants.Edges)-1].Cursor
		nextPageData, err := s.getPage(ctx, id, cursor)
		if err != nil {
			return nil, fmt.Errorf("get page: %w", err)
		}
		out.Variants.Edges = append(out.Variants.Edges, nextPageData.Variants.Edges...)
		hasNextPage = nextPageData.Variants.PageInfo.HasNextPage
	}

	return out, nil
}

func (s *ProductServiceOp) getPage(ctx context.Context, id string, cursor string) (*model.Product, error) {
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
	err := s.client.gql.QueryString(ctx, q, vars, &out)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out.Product, nil
}

func (s *ProductServiceOp) Create(ctx context.Context, product model.ProductCreateInput, media []model.CreateMediaInput) (*string, error) {
	m := mutationProductCreate{}

	vars := map[string]interface{}{
		"product": product,
		"media":   media,
	}

	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductCreateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.ProductCreateResult.UserErrors)
	}

	return &m.ProductCreateResult.Product.ID, nil
}

func (s *ProductServiceOp) Update(ctx context.Context, product model.ProductUpdateInput, media []model.CreateMediaInput) error {
	m := mutationProductUpdate{}

	vars := map[string]interface{}{
		"product": product,
		"media":   media,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductUpdateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) Delete(ctx context.Context, product model.ProductDeleteInput) error {
	m := mutationProductDelete{}

	vars := map[string]interface{}{
		"input": product,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductDeleteResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductDeleteResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) VariantsBulkCreate(ctx context.Context, id string, input []model.ProductVariantsBulkInput, strategy model.ProductVariantsBulkCreateStrategy) error {
	m := mutationProductVariantsBulkCreate{}

	vars := map[string]interface{}{
		"productId": id,
		"variants":  input,
		"strategy":  strategy,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductVariantsBulkCreateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantsBulkCreateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) VariantsBulkUpdate(ctx context.Context, id string, input []model.ProductVariantsBulkInput) error {
	m := mutationProductVariantsBulkUpdate{}

	vars := map[string]interface{}{
		"productId": id,
		"variants":  input,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductVariantsBulkUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantsBulkUpdateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) VariantsBulkReorder(ctx context.Context, id string, input []model.ProductVariantPositionInput) error {
	m := mutationProductVariantsBulkReorder{}

	vars := map[string]interface{}{
		"productId": id,
		"positions": input,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.ProductVariantsBulkReorderResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantsBulkReorderResult.UserErrors)
	}

	return nil
}
