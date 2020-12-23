package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/graphql"
)

type InventoryService interface {
	Update(id graphql.ID, input InventoryItemUpdateInput) error
	Adjust(locationID graphql.ID, input []InventoryAdjustItemInput) error
	ActivateInventory(locationID graphql.ID, id graphql.ID) error
}

type InventoryServiceOp struct {
	client *Client
}

type InventoryItem struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
	SKU              graphql.String `json:"sku,omitempty"`
}

type InventoryLevel struct {
	UpdatedAt graphql.String `json:"updatedAt,omitempty"`
	Available graphql.Int    `json:"available,omitempty"`
	Item      InventoryItem  `json:"item,omitempty"`
}

type InventoryItemUpdateInput struct {
	Cost graphql.Float `json:"cost,omitempty"`
}

type mutationInventoryItemUpdate struct {
	InventoryItemUpdateResult InventoryItemUpdateResult `graphql:"inventoryItemUpdate(id: $id, input: $input)" json:"inventoryItemUpdate"`
}

type InventoryItemUpdateResult struct {
	UserErrors []UserErrors `json:"userErrors,omitempty"`
}

type InventoryAdjustItemInput struct {
	InventoryItemID graphql.ID  `json:"inventoryItemId,omitempty"`
	AvailableDelta  graphql.Int `json:"availableDelta,omitempty"`
}

type mutationInventoryBulkAdjustQuantityAtLocation struct {
	InventoryBulkAdjustQuantityAtLocationResult InventoryBulkAdjustQuantityAtLocationResult `graphql:"inventoryBulkAdjustQuantityAtLocation(locationId: $locationId, inventoryItemAdjustments: $inventoryItemAdjustments)" json:"inventoryBulkAdjustQuantityAtLocation"`
}

type InventoryBulkAdjustQuantityAtLocationResult struct {
	UserErrors []UserErrors `json:"userErrors,omitempty"`
}

type mutationInventoryActivate struct {
	InventoryActivateResult InventoryActivateResult `graphql:"inventoryActivate(inventoryItemId: $itemID, locationId: $locationId)" json:"inventoryActivate"`
}

type InventoryActivateResult struct {
	UserErrors []UserErrors `json:"userErrors,omitempty"`
}

func (s *InventoryServiceOp) Update(id graphql.ID, input InventoryItemUpdateInput) error {
	m := mutationInventoryItemUpdate{}
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.InventoryItemUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryItemUpdateResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) Adjust(locationID graphql.ID, input []InventoryAdjustItemInput) error {
	m := mutationInventoryBulkAdjustQuantityAtLocation{}
	vars := map[string]interface{}{
		"locationId":               locationID,
		"inventoryItemAdjustments": input,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) ActivateInventory(locationID graphql.ID, id graphql.ID) error {
	m := mutationInventoryActivate{}
	vars := map[string]interface{}{
		"itemID":     id,
		"locationId": locationID,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.InventoryActivateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryActivateResult.UserErrors)
	}

	return nil
}
