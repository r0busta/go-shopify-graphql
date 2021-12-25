package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/v2/graph/model"
)

//go:generate mockgen -destination=./mock/inventory_service.go -package=mock . InventoryService
type InventoryService interface {
	Update(id string, input model.InventoryItemUpdateInput) error
	Adjust(locationID string, input []model.InventoryAdjustItemInput) error
	ActivateInventory(locationID string, id string) error
}

type InventoryServiceOp struct {
	client *Client
}

var _ InventoryService = &InventoryServiceOp{}

type mutationInventoryItemUpdate struct {
	InventoryItemUpdateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryItemUpdate(id: $id, input: $input)" json:"inventoryItemUpdate"`
}

type mutationInventoryBulkAdjustQuantityAtLocation struct {
	InventoryBulkAdjustQuantityAtLocationResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryBulkAdjustQuantityAtLocation(locationId: $locationId, inventoryItemAdjustments: $inventoryItemAdjustments)" json:"inventoryBulkAdjustQuantityAtLocation"`
}

type mutationInventoryActivate struct {
	InventoryActivateResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryActivate(inventoryItemId: $itemID, locationId: $locationId)" json:"inventoryActivate"`
}

func (s *InventoryServiceOp) Update(id string, input model.InventoryItemUpdateInput) error {
	m := mutationInventoryItemUpdate{}
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryItemUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryItemUpdateResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) Adjust(locationID string, input []model.InventoryAdjustItemInput) error {
	m := mutationInventoryBulkAdjustQuantityAtLocation{}
	vars := map[string]interface{}{
		"locationId":               locationID,
		"inventoryItemAdjustments": input,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) ActivateInventory(locationID string, id string) error {
	m := mutationInventoryActivate{}
	vars := map[string]interface{}{
		"itemID":     id,
		"locationId": locationID,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryActivateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryActivateResult.UserErrors)
	}

	return nil
}
