package shopify

import (
	"context"
	"fmt"
	"log"

	"github.com/shurcooL/graphql"
)

type CollectionService interface {
	ListAll() ([]*Collection, error)

	Get(id graphql.ID) (*Collection, error)

	Create(collection *CollectionCreate) (graphql.ID, error)
	CreateBulk(collections []*CollectionCreate) error

	Update(collection *CollectionCreate) error
}

type CollectionServiceOp struct {
	client *Client
}

type Collection struct {
	ID            graphql.ID     `json:"id,omitempty"`
	Handle        graphql.String `json:"handle,omitempty"`
	ProductsCount graphql.Int    `json:"productsCount,omitempty"`

	Products struct {
		Edges []ProductShortNode `json:"edges,omitempty"`
	} `graphql:"products(first: 100)" json:"products,omitempty"`
}

type ProductShortNode struct {
	Node ProductShort `json:"node,omitempty"`
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

func (s *CollectionServiceOp) ListAll() ([]*Collection, error) {
	query := `
		{
			collections{
				edges{
					node{
						id
						handle
					}
				}
			}
		}
`

	res := []*Collection{}
	err := s.client.BulkOperation.BulkQuery(query, &res)
	if err != nil {
		return []*Collection{}, err
	}

	return res, nil
}

func (s *CollectionServiceOp) Get(id graphql.ID) (*Collection, error) {
	var q struct {
		Collection `graphql:"collection(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": id,
	}

	err := s.client.gql.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, err
	}

	return &q.Collection, nil
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
