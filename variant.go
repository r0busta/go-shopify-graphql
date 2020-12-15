package shopify

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

type VariantService interface {
	Update(variant *ProductVariantUpdate) error
}

type VariantServiceOp struct {
	client *Client
}

type ProductVariant struct {
	ID                graphql.ID       `json:"id,omitempty"`
	LegacyResourceID  graphql.String   `json:"legacyResourceId,omitempty"`
	SKU               graphql.String   `json:"sku,omitempty"`
	SelectedOptions   []SelectedOption `json:"selectedOptions,omitempty"`
	CompareAtPrice    Money            `json:"compareAtPrice,omitempty"`
	Price             Money            `json:"price,omitempty"`
	InventoryQuantity graphql.Int      `json:"inventoryQuantity,omitempty"`
	InventoryItem     InventoryItem    `json:"inventoryItem,omitempty"`
}

type SelectedOption struct {
	Name  graphql.String `json:"name,omitempty"`
	Value graphql.String `json:"value,omitempty"`
}

type ProductVariantPricePair struct {
	CompareAtPrice Money `json:"compareAtPrice,omitempty"`
	Price          Money `json:"price,omitempty"`
}

type ProductVariantUpdate struct {
	ProductVariantInput ProductVariantInput
}

type ProductVariantInput struct {
	// The value of the barcode associated with the product.
	Barcode graphql.String `json:"barcode,omitempty"`

	// The compare-at price of the variant.
	CompareAtPrice *Money `json:"compareAtPrice"`

	// The ID of the fulfillment service associated with the variant.
	FulfillmentServiceID graphql.ID `json:"fulfillmentServiceId,omitempty"`

	// The Harmonized System Code (or HS Tariff Code) for the variant.
	HarmonizedSystemCode graphql.String `json:"harmonizedSystemCode,omitempty"`

	// Specifies the product variant to update or create a new variant if absent.
	ID graphql.ID `json:"id,omitempty"`

	// The ID of the image that's associated with the variant.
	ImageID graphql.ID `json:"imageId,omitempty"`

	// The URL of an image to associate with the variant. This field can only be used through mutations that create product images and must match one of the URLs being created on the product.
	ImageSrc graphql.String `json:"imageSrc,omitempty"`

	// Inventory Item associated with the variant, used for unit cost.
	InventoryItem *InventoryItemInput `json:"inventoryItem,omitempty"`

	// Whether customers are allowed to place an order for the product variant when it's out of stock.
	InventoryPolicy ProductVariantInventoryPolicy `json:"inventoryPolicy,omitempty"`

	// Create only field. The inventory quantities at each location where the variant is stocked.
	InventoryQuantities []InventoryLevelInput `json:"inventoryQuantities,omitempty"`

	// The ID of the corresponding resource in the REST Admin API.
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`

	// Additional customizable information about the product variant.
	Metafields []MetafieldInput `json:"metafields,omitempty"`

	// The custom properties that a shop owner uses to define product variants.
	Options []graphql.String `json:"options,omitempty"`

	// The order of the product variant in the list of product variants. The first position in the list is 1.
	Position graphql.Int `json:"position,omitempty"`

	// The price of the variant.
	Price Money `json:"price,omitempty"`

	// Create only required field. Specifies the product on which to create the variant.
	ProductID graphql.ID `json:"productId,omitempty"`

	// Whether the variant requires shipping.
	RequiresShipping graphql.Boolean `json:"requiresShipping,omitempty"`

	// The SKU for the variant.
	SKU graphql.String `json:"sku,omitempty"`

	// This parameter applies only to the stores that have the Avalara AvaTax app installed. Specifies the Avalara tax code for the product variant.
	TaxCode graphql.String `json:"taxCode,omitempty"`

	// Whether the variant is taxable.
	Taxable graphql.Boolean `json:"taxable,omitempty"`

	// This argument is deprecated: Variant title is not a writable field; it is generated from the selected variant options.
	Title graphql.String `json:"title,omitempty"`

	// The weight of the variant.
	Weight graphql.Float `json:"weight,omitempty"`

	// The unit of weight that's used to measure the variant.
	WeightUnit WeightUnit `json:"weightUnit,omitempty"`
}

type InventoryItemInput struct {
	// Unit cost associated with the inventory item, the currency is the shop's default currency.
	Cost Decimal `json:"cost,omitempty"`
	// Whether the inventory item is tracked. If true, then inventory quantity changes are tracked by Shopify.
	Tracked graphql.Boolean `json:"tracked,omitempty"`
}

// ProductVariantInventoryPolicy String enum: CONTINUE, DENY
type ProductVariantInventoryPolicy string

type InventoryLevelInput struct {
	AvailableQuantity graphql.Int `json:"availableQuantity"`
	LocationID        graphql.ID  `json:"locationId"`
}

// WeightUnit String enum: GRAMS, KILOGRAMS, OUNCES, POUNDS
type WeightUnit string

type mutationProductVariantUpdate struct {
	ProductVariantUpdateResult productVariantUpdateResult `graphql:"productVariantUpdate(input: $input)"`
}

type productVariantUpdateResult struct {
	UserErrors []UserErrors
}

func (s *VariantServiceOp) Update(variant *ProductVariantUpdate) error {
	m := mutationProductVariantUpdate{}

	vars := map[string]interface{}{
		"input": variant.ProductVariantInput,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.ProductVariantUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantUpdateResult.UserErrors)
	}

	return nil
}
