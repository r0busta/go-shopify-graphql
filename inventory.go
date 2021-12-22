package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	"github.com/r0busta/graphql"
)

//go:generate mockgen -destination=./mock/inventory_service.go -package=mock . InventoryService
type InventoryService interface {
	Update(id graphql.ID, input model.InventoryItemUpdateInput) error
	Adjust(locationID graphql.ID, input []model.InventoryAdjustItemInput) error
	ActivateInventory(locationID graphql.ID, id graphql.ID) error
}

type InventoryServiceOp struct {
	client *Client
}

type mutationInventoryItemUpdate struct {
	InventoryItemUpdateResult model.InventoryItemUpdatePayload `graphql:"inventoryItemUpdate(id: $id, input: $input)" json:"inventoryItemUpdate"`
}

type mutationInventoryBulkAdjustQuantityAtLocation struct {
	InventoryBulkAdjustQuantityAtLocationResult model.InventoryBulkAdjustQuantityAtLocationPayload `graphql:"inventoryBulkAdjustQuantityAtLocation(locationId: $locationId, inventoryItemAdjustments: $inventoryItemAdjustments)" json:"inventoryBulkAdjustQuantityAtLocation"`
}

type mutationInventoryActivate struct {
	InventoryActivateResult model.InventoryActivatePayload `graphql:"inventoryActivate(inventoryItemId: $itemID, locationId: $locationId)" json:"inventoryActivate"`
}

func (s *InventoryServiceOp) Update(id graphql.ID, input model.InventoryItemUpdateInput) error {
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

func (s *InventoryServiceOp) Adjust(locationID graphql.ID, input []model.InventoryAdjustItemInput) error {
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
