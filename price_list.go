//go:generate mockgen -package shopify -destination price_list_mock.go -source price_list.go

package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
)

// PriceList represents a price list.
type PriceList struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

// PriceListService defines the price list service operations.
type PriceListService interface {
	GetAll(ctx context.Context) ([]PriceList, error)
	AddFixedPrices(ctx context.Context, priceListID string, prices []PriceListPriceInput) error
}

// PriceListBulkQueryClient defines the required bulk query client operations.
type PriceListBulkQueryClient interface {
	BulkQuery(query string, v interface{}) error
}

// PriceListMutationClient defines the required mutation client operations.
type PriceListMutationClient interface {
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error
}

// PriceListServiceOp represents a price list service.
type PriceListServiceOp struct {
	bulkQueryClient PriceListBulkQueryClient
	mutationClient  PriceListMutationClient
}

const (
	priceListsGetAllBulkQuery = `
		{
			priceLists {
				edges {
      				node {
        				id
        				name
        				currency
      				}
    			}
  			}
		}
	`
)

// GetAll returns all price lists for the store.
func (s *PriceListServiceOp) GetAll(_ context.Context) ([]PriceList, error) {
	var res []PriceList
	err := s.bulkQueryClient.BulkQuery(priceListsGetAllBulkQuery, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve price lists: %w", err)
	}

	return res, nil
}

// PriceListPriceInput represents a price list price input.
type PriceListPriceInput struct {
	CompareAtPrice model.MoneyInput `json:"compareAtPrice"`
	Price          model.MoneyInput `json:"price"`
	VariantID      string           `json:"variantId"`
}

type priceListPriceOutput struct {
	CompareAtPrice model.MoneyV2 `json:"compareAtPrice"`
	Price          model.MoneyV2 `json:"price"`
	// Commented because I'm not sure why this causes a deadlock on the query generator, and we don't need it now...
	// Variant model.ProductVariant `json:"variant"`
}

type priceListFixedPricesAddResult struct {
	Prices     []priceListPriceOutput `json:"prices"`
	UserErrors []model.UserError      `json:"userErrors,omitempty"`
}

type mutationPriceListFixedPricesAdd struct {
	PriceListFixedPricesAddResult priceListFixedPricesAddResult `graphql:"priceListFixedPricesAdd(priceListId: $priceListId, prices: $prices)" json:"priceListFixedPricesAdd"`
}

// AddFixedPrices adds a fixed prices to a price list based on the arguments received.
func (s *PriceListServiceOp) AddFixedPrices(ctx context.Context, priceListID string, prices []PriceListPriceInput) error {
	m := mutationPriceListFixedPricesAdd{}

	vars := map[string]interface{}{
		"priceListId": priceListID,
		"prices":      prices,
	}
	err := s.mutationClient.Mutate(ctx, &m, vars)
	if err != nil {
		return err
	}

	if len(m.PriceListFixedPricesAddResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.PriceListFixedPricesAddResult.UserErrors)
	}

	return nil
}
