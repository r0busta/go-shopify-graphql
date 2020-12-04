package shopify

import (
	"strings"

	"github.com/shurcooL/graphql"
)

type OrderService interface {
	List(query string) ([]*Order, error)
	ListAll() ([]*Order, error)
}

type OrderServiceOp struct {
	client *Client
}

type Order struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
	CreatedAt        DateTime       `json:"createdAt,omitempty"`
	Customer         Customer       `json:"customer,omitempty"`
	ClientIP         graphql.String `json:"clientIp,omitempty"`
	TaxLines         []TaxLine      `json:"taxLines,omitempty"`
	TotalReceivedSet MoneyBag       `json:"totalReceivedSet,omitempty"`
	LineItems        []LineItem     `json:"lineItems,omitempty"`
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
	Quantity               graphql.Int     `json:"quantity,omitempty"`
	Product                LineItemProduct `json:"product,omitempty"`
	Variant                LineItemVariant `json:"variant,omitempty"`
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

func (s *OrderServiceOp) List(query string) ([]*Order, error) {
	q := `
		{
			orders(query: "$query"){
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
