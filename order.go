package shopify

import (
	"context"
	"fmt"
	"strings"

	"github.com/r0busta/graphql"
)

type OrderService interface {
	List(query string) ([]*Order, error)
	ListAll() ([]*Order, error)

	Update(input *OrderInput) error

	GetFulfillmentOrdersAtLocation(orderID graphql.ID, locationID graphql.ID) ([]FulfillmentOrder, error)
}

type OrderServiceOp struct {
	client *Client
}

type Order struct {
	ID                graphql.ID         `json:"id,omitempty"`
	LegacyResourceID  graphql.String     `json:"legacyResourceId,omitempty"`
	Name              graphql.String     `json:"name,omitempty"`
	CreatedAt         DateTime           `json:"createdAt,omitempty"`
	Customer          Customer           `json:"customer,omitempty"`
	ClientIP          graphql.String     `json:"clientIp,omitempty"`
	TaxLines          []TaxLine          `json:"taxLines,omitempty"`
	TotalReceivedSet  MoneyBag           `json:"totalReceivedSet,omitempty"`
	LineItems         []LineItem         `json:"lineItems,omitempty"`
	ShippingAddress   MailingAddress     `json:"shippingAddress,omitempty"`
	ShippingLine      ShippingLine       `json:"shippingLine,omitempty"`
	Note              graphql.String     `json:"note,omitempty"`
	Tags              []graphql.String   `json:"tags,omitempty"`
	FulfillmentOrders []FulfillmentOrder `json:"fulfillmentOrders,omitempty"`
}

type ShippingLine struct {
	Title            graphql.String `json:"title,omitempty"`
	OriginalPriceSet MoneyBag       `json:"originalPriceSet,omitempty"`
}

type TaxLine struct {
	PriceSet       MoneyBag       `json:"priceSet,omitempty"`
	Rate           graphql.Float  `json:"rate,omitempty"`
	RatePercentage graphql.Float  `json:"ratePercentage,omitempty"`
	Title          graphql.String `json:"title,omitempty"`
}

type OrderLineItemNode struct {
	Node LineItem `json:"node,omitempty"`
}

type LineItem struct {
	ID                     graphql.ID      `json:"id,omitempty"`
	SKU                    graphql.String  `json:"sku,omitempty"`
	Quantity               graphql.Int     `json:"quantity,omitempty"`
	FulfillableQuantity    graphql.Int     `json:"fulfillableQuantity,omitempty"`
	Vendor                 graphql.String  `json:"vendor,omitempty"`
	Title                  graphql.String  `json:"title,omitempty"`
	VariantTitle           graphql.String  `json:"variantTitle,omitempty"`
	Product                LineItemProduct `json:"product,omitempty"`
	Variant                LineItemVariant `json:"variant,omitempty"`
	OriginalTotalSet       MoneyBag        `json:"originalTotalSet,omitempty"`
	OriginalUnitPriceSet   MoneyBag        `json:"originalUnitPriceSet,omitempty"`
	DiscountedUnitPriceSet MoneyBag        `json:"discountedUnitPriceSet,omitempty"`
}

type LineItemProduct struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
}

type LineItemVariant struct {
	ID               graphql.ID       `json:"id,omitempty"`
	LegacyResourceID graphql.String   `json:"legacyResourceId,omitempty"`
	SelectedOptions  []SelectedOption `json:"selectedOptions,omitempty"`
}

type FulfillmentOrder struct {
	ID                        graphql.ID                 `json:"id,omitempty"`
	Status                    FulfillmentOrderStatus     `json:"status,omitempty"`
	FulfillmentOrderLineItems []FulfillmentOrderLineItem `json:"lineItems,omitempty"`
}

type FulfillmentOrderStatus string

type FulfillmentOrderLineItem struct {
	ID                graphql.ID  `json:"id,omitempty"`
	RemainingQuantity graphql.Int `json:"remainingQuantity,omitempty"`
	LineItem          LineItem    `json:"lineItem,omitempty"`
}

type mutationOrderUpdate struct {
	OrderUpdateResult OrderUpdateResult `graphql:"orderUpdate(input: $input)" json:"orderUpdate"`
}

type OrderUpdateResult struct {
	UserErrors []UserErrors `json:"userErrors"`
}
type OrderInput struct {
	ID   graphql.ID       `json:"id,omitempty"`
	Tags []graphql.String `json:"tags,omitempty"`
}

func (s *OrderServiceOp) List(query string) ([]*Order, error) {
	q := `
		{
			orders(query: "$query"){
				edges{
					node{
						id
						legacyResourceId
						name
						createdAt
						customer{
							id
							legacyResourceId
							firstName
							displayName
							email
						}
						clientIp
						shippingAddress{
							address1
							address2
							city
							province
							country
							zip
						}
						note
						shippingLine{
							originalPriceSet{
								presentmentMoney{
									amount
									currencyCode
								}
								shopMoney{
									amount
									currencyCode
								}
							}
							title
						}
						taxLines{
							priceSet{
								presentmentMoney{
									amount
									currencyCode
								}
								shopMoney{
									amount
									currencyCode
								}
							}
							rate
							ratePercentage
							title
						}
						totalReceivedSet{
							presentmentMoney{
								amount
								currencyCode
							}
							shopMoney{
								amount
								currencyCode
							}
						}
						lineItems{
							edges{
								node{
									id
									sku
									quantity
									fulfillableQuantity
									product{
										id
										legacyResourceId										
									}
									vendor
									title
									variantTitle
									variant{
										id
										legacyResourceId	
										selectedOptions{
											name
											value
										}									
									}
									originalTotalSet{
										presentmentMoney{
											amount
											currencyCode
										}
										shopMoney{
											amount
											currencyCode
										}
									}
									originalUnitPriceSet{
										presentmentMoney{
											amount
											currencyCode
										}
										shopMoney{
											amount
											currencyCode
										}
									}
									discountedUnitPriceSet{
										presentmentMoney{
											amount
											currencyCode
										}
										shopMoney{
											amount
											currencyCode
										}
									}
								}
							}
						}
					}
				}
			}
		}
`
	q = strings.ReplaceAll(q, "$query", query)

	res := []*Order{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*Order{}, err
	}

	return res, nil
}

func (s *OrderServiceOp) ListAll() ([]*Order, error) {
	q := `
		{
			orders{
				edges{
					node{
						id
						legacyResourceId
						createdAt
						customer{
							id
							legacyResourceId
						}
						clientIp
						taxLines{
							priceSet{
								presentmentMoney{
									amount
									currencyCode
								}
								shopMoney{
									amount
									currencyCode
								}
							}
							rate
							ratePercentage
							title
						}
						totalReceivedSet{
							presentmentMoney{
								amount
								currencyCode
							}
							shopMoney{
								amount
								currencyCode
							}
						}
						lineItems{
							edges{
								node{
									id
									quantity
									product{
										id
										legacyResourceId										
									}
									variant{
										id
										legacyResourceId	
										selectedOptions{
											name
											value
										}									
									}
									originalUnitPriceSet{
										presentmentMoney{
											amount
											currencyCode
										}
										shopMoney{
											amount
											currencyCode
										}
									}
									discountedUnitPriceSet{
										presentmentMoney{
											amount
											currencyCode
										}
										shopMoney{
											amount
											currencyCode
										}
									}
								}
							}
						}
					}
				}
			}
		}
`

	res := []*Order{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []*Order{}, err
	}

	return res, nil
}

func (s *OrderServiceOp) Update(input *OrderInput) error {
	m := mutationOrderUpdate{}

	vars := map[string]interface{}{
		"input": input,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.OrderUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.OrderUpdateResult.UserErrors)
	}

	return nil
}

func (s *OrderServiceOp) GetFulfillmentOrdersAtLocation(orderID graphql.ID, locationID graphql.ID) ([]FulfillmentOrder, error) {
	q := `
	{
		order(id:"$id"){
			fulfillmentOrders(query:"$query"){
				edges {
					node {
						id
						status
						lineItems{
							edges {
								node {
									id
									remainingQuantity
									lineItem{
										sku
									}								
								}
							}
						}
					}
				}
			}
		}
	}`

	q = strings.ReplaceAll(q, "$id", orderID.(string))
	q = strings.ReplaceAll(q, "$query", fmt.Sprintf(`assigned_location_id:%s`, locationID.(string)))
	res := []FulfillmentOrder{}
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return []FulfillmentOrder{}, err
	}

	return res, nil
}
