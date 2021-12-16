package shopify

import (
	"fmt"
)

// PriceList represents a price list.
type PriceList struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

// PriceListService defines the price list service operations.
type PriceListService interface {
	GetAll() ([]PriceList, error)
}

// PriceListServiceOp represents a price list service.
type PriceListServiceOp struct {
	service BulkOperationService
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
func (s *PriceListServiceOp) GetAll() ([]PriceList, error) {
	var res []PriceList
	err := s.service.BulkQuery(priceListsGetAllBulkQuery, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve price lists: %w", err)
	}

	return res, nil
}

// AddFixedPrices TODO
func (s *PriceListServiceOp) AddFixedPrices() error {
	// TODO
	return nil
}
