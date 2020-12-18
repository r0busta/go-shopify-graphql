package shopify

import (
	"context"
	"fmt"
	"log"

	"github.com/r0busta/graphql"
)

type CollectionService interface {
	ListAll() ([]*CollectionBulkResult, error)

	Get(id graphql.ID) (*CollectionQueryResult, error)

	Create(collection *CollectionCreate) (graphql.ID, error)
	CreateBulk(collections []*CollectionCreate) error

	Update(collection *CollectionCreate) error
}

type CollectionServiceOp struct {
	client *Client
}

type CollectionBase struct {
	ID            graphql.ID     `json:"id,omitempty"`
	Handle        graphql.String `json:"handle,omitempty"`
	Title         graphql.String `json:"title,omitempty"`
	ProductsCount graphql.Int    `json:"productsCount,omitempty"`
}

type CollectionBulkResult struct {
	CollectionBase

	Products []ProductBulkResult `json:"products,omitempty"`
}

type CollectionQueryResult struct {
	CollectionBase

	Products struct {
		Edges []struct {
			Product ProductQueryResult `json:"node,omitempty"`
			Cursor  string             `json:"cursor,omitempty"`
		} `json:"edges,omitempty"`
		PageInfo PageInfo `json:"pageInfo,omitempty"`
	} `json:"products,omitempty"`
}

type CollectionCreate struct {
	CollectionInput CollectionInput
}

type mutationCollectionCreate struct {
	CollectionCreateResult CollectionCreateResult `graphql:"collectionCreate(input: $input)"`
}

type mutationCollectionUpdate struct {
	CollectionCreateResult CollectionCreateResult `graphql:"collectionUpdate(input: $input)"`
}

type CollectionInput struct {
	// The description of the collection, in HTML format.
	DescriptionHTML graphql.String `json:"descriptionHtml,omitempty"`

	// A unique human-friendly string for the collection. Automatically generated from the collection's title.
	Handle graphql.String `json:"handle,omitempty"`

	// Specifies the collection to update or create a new collection if absent.
	ID graphql.ID `json:"id,omitempty"`

	// The image associated with the collection.
	Image *ImageInput `json:"image,omitempty"`

	// The metafields to associate with this collection.
	Metafields []MetafieldInput `json:"metafields,omitempty"`

	// Initial list of collection products. Only valid with productCreate and without rules.
	Products []graphql.ID `json:"products,omitempty"`

	// Indicates whether a redirect is required after a new handle has been provided. If true, then the old handle is redirected to the new one automatically.
	RedirectNewHandle graphql.Boolean `json:"redirectNewHandle,omitempty"`

	//	The rules used to assign products to the collection.
	RuleSet *CollectionRuleSetInput `json:"ruleSet,omitempty"`

	// SEO information for the collection.
	SEO *SEOInput `json:"seo,omitempty"`

	// The order in which the collection's products are sorted.
	SortOrder *CollectionSortOrder `json:"sortOrder,omitempty"`

	// The theme template used when viewing the collection in a store.
	TemplateSuffix graphql.String `json:"templateSuffix,omitempty"`

	// Required for creating a new collection.
	Title graphql.String `json:"title,omitempty"`
}

type CollectionRuleSetInput struct {
	// Whether products must match any or all of the rules to be included in the collection. If true, then products must match one or more of the rules to be included in the collection. If false, then products must match all of the rules to be included in the collection.
	AppliedDisjunctively graphql.Boolean `json:"appliedDisjunctively"` // REQUIRED

	// The rules used to assign products to the collection.
	Rules []CollectionRuleInput `json:"rules,omitempty"`
}

type CollectionRuleInput struct {
	// The attribute that the rule focuses on (for example, title or product_type).
	Column CollectionRuleColumn `json:"column,omitempty"` // REQUIRED

	// The value that the operator is applied to (for example, Hats).
	Condition graphql.String `json:"condition,omitempty"` // REQUIRED

	// The type of operator that the rule is based on (for example, equals, contains, or not_equals).
	Relation CollectionRuleRelation `json:"relation,omitempty"` // REQUIRED
}

// CollectionRuleColumn string enum
// VENDOR The vendor attribute.
// TAG The tag attribute.
// TITLE The title attribute.
// TYPE The type attribute.
// VARIANT_COMPARE_AT_PRICE The variant_compare_at_price attribute.
// VARIANT_INVENTORY The variant_inventory attribute.
// VARIANT_PRICE The variant_price attribute.
// VARIANT_TITLE The variant_title attribute.
// VARIANT_WEIGHT The variant_weight attribute.
// IS_PRICE_REDUCED The is_price_reduced attribute.
type CollectionRuleColumn string

// CollectionRuleRelation enum
// STARTS_WITH The attribute starts with the condition.
// ENDS_WITH The attribute ends with the condition.
// EQUALS The attribute is equal to the condition.
// GREATER_THAN The attribute is greater than the condition.
// IS_NOT_SET The attribute is not set.
// IS_SET The attribute is set.
// LESS_THAN The attribute is less than the condition.
// NOT_CONTAINS The attribute does not contain the condition.
// NOT_EQUALS The attribute does not equal the condition.
// CONTAINS The attribute contains the condition.
type CollectionRuleRelation string

// CollectionSortOrder enum
// PRICE_DESC By price, in descending order (highest - lowest).
// ALPHA_DESC Alphabetically, in descending order (Z - A).
// BEST_SELLING By best-selling products.
// CREATED By date created, in ascending order (oldest - newest).
// CREATED_DESC By date created, in descending order (newest - oldest).
// MANUAL In the order set manually by the merchant.
// PRICE_ASC By price, in ascending order (lowest - highest).
// ALPHA_ASC Alphabetically, in ascending order (A - Z).
type CollectionSortOrder string

type CollectionCreateResult struct {
	Collection struct {
		ID graphql.ID `json:"id,omitempty"`
	}
	UserErrors []UserErrors
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

func (s *CollectionServiceOp) ListAll() ([]*CollectionBulkResult, error) {
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

	res := []*CollectionBulkResult{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*CollectionBulkResult{}, err
	}

	return res, nil
}

func (s *CollectionServiceOp) Get(id graphql.ID) (*CollectionQueryResult, error) {
	out, err := s.getPage(id, "")
	if err != nil {
		return nil, err
	}

	nextPageData := out
	hasNextPage := out.Products.PageInfo.HasNextPage
	for hasNextPage && len(nextPageData.Products.Edges) > 0 {
		cursor := nextPageData.Products.Edges[len(nextPageData.Products.Edges)-1].Cursor
		nextPageData, err := s.getPage(id, cursor)
		if err != nil {
			return nil, err
		}
		out.Products.Edges = append(out.Products.Edges, nextPageData.Products.Edges...)
		hasNextPage = nextPageData.Products.PageInfo.HasNextPage
	}

	return out, nil
}

func (s *CollectionServiceOp) getPage(id graphql.ID, cursor string) (*CollectionQueryResult, error) {
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
		Collection *CollectionQueryResult `json:"collection"`
	}{}
	err := s.client.gql.QueryString(context.Background(), q, vars, &out)
	if err != nil {
		return nil, err
	}

	return out.Collection, nil
}

func (s *CollectionServiceOp) CreateBulk(collections []*CollectionCreate) error {
	for _, c := range collections {
		_, err := s.client.Collection.Create(c)
		if err != nil {
			log.Printf("Warning! Couldn't create collection (%v): %s", c, err)
		}
	}

	return nil
}

func (s *CollectionServiceOp) Create(collection *CollectionCreate) (graphql.ID, error) {
	var id graphql.ID
	m := mutationCollectionCreate{}

	vars := map[string]interface{}{
		"input": collection.CollectionInput,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return id, err
	}

	if len(m.CollectionCreateResult.UserErrors) > 0 {
		return id, fmt.Errorf("%+v", m.CollectionCreateResult.UserErrors)
	}

	id = m.CollectionCreateResult.Collection.ID
	return id, nil
}

func (s *CollectionServiceOp) Update(collection *CollectionCreate) error {
	m := mutationCollectionUpdate{}

	vars := map[string]interface{}{
		"input": collection.CollectionInput,
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
