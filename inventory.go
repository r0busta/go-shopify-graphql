package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/v3/graph/model"
)

//go:generate mockgen -destination=./mock/inventory_service.go -package=mock . InventoryService
type InventoryService interface {
	Update(ctx context.Context, id string, input model.InventoryItemUpdateInput) error
	Adjust(ctx context.Context, locationID string, input []model.InventoryAdjustItemInput) error
	AdjustQuantities(ctx context.Context, reason, name string, referenceDocumentUri *string, changes []model.InventoryChangeInput) error
	SetOnHandQuantities(ctx context.Context, reason string, referenceDocumentUri *string, setQuantities []model.InventorySetQuantityInput) error
	ActivateInventory(ctx context.Context, locationID string, id string) error
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

type mutationInventoryAdjustQuantities struct {
	InventoryAdjustQuantitiesResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventoryAdjustQuantities(input: $input)" json:"inventoryAdjustQuantities"`
}

type mutationInventorySetOnHandQuantities struct {
	InventorySetOnHandQuantitiesResult struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"inventorySetOnHandQuantities(input: $input)" json:"inventorySetOnHandQuantities"`
}

func (s *InventoryServiceOp) Update(ctx context.Context, id string, input model.InventoryItemUpdateInput) error {
	m := mutationInventoryItemUpdate{}
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryItemUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryItemUpdateResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) Adjust(ctx context.Context, locationID string, input []model.InventoryAdjustItemInput) error {
	m := mutationInventoryBulkAdjustQuantityAtLocation{}
	vars := map[string]interface{}{
		"locationId":               locationID,
		"inventoryItemAdjustments": input,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryBulkAdjustQuantityAtLocationResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) AdjustQuantities(ctx context.Context, reason, name string, referenceDocumentUri *string, changes []model.InventoryChangeInput) error {
	m := mutationInventoryAdjustQuantities{}
	vars := map[string]interface{}{
		"input": model.InventoryAdjustQuantitiesInput{
			Name:                 name,
			Reason:               reason,
			ReferenceDocumentURI: referenceDocumentUri,
			Changes:              changes,
		},
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryAdjustQuantitiesResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryAdjustQuantitiesResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) SetOnHandQuantities(ctx context.Context, reason string, referenceDocumentUri *string, setQuantities []model.InventorySetQuantityInput) error {
	m := mutationInventorySetOnHandQuantities{}
	vars := map[string]interface{}{
		"input": model.InventorySetOnHandQuantitiesInput{
			Reason:               reason,
			ReferenceDocumentURI: referenceDocumentUri,
			SetQuantities:        setQuantities,
		},
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventorySetOnHandQuantitiesResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventorySetOnHandQuantitiesResult.UserErrors)
	}

	return nil
}

func (s *InventoryServiceOp) ActivateInventory(ctx context.Context, locationID string, id string) error {
	m := mutationInventoryActivate{}
	vars := map[string]interface{}{
		"itemID":     id,
		"locationId": locationID,
	}
	err := s.client.gql.Mutate(ctx, &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %w", err)
	}

	if len(m.InventoryActivateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.InventoryActivateResult.UserErrors)
	}

	return nil
}
