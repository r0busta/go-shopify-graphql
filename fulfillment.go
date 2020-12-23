package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/graphql"
)

type FulfillmentService interface {
	Create(input FulfillmentV2Input) error
}

type FulfillmentServiceOp struct {
	client *Client
}

type FulfillmentV2Input struct {
	LineItemsByFulfillmentOrder []FulfillmentOrderLineItemsInput `json:"lineItemsByFulfillmentOrder,omitempty"`
	NotifyCustomer              graphql.Boolean                  `json:"notifyCustomer,omitempty"`
	TrackingInfo                FulfillmentTrackingInput         `json:"trackingInfo,omitempty"`
}

type FulfillmentOrderLineItemsInput struct {
	FulfillmentOrderID        graphql.ID                      `json:"fulfillmentOrderId,omitempty"`
	FulfillmentOrderLineItems []FulfillmentOrderLineItemInput `json:"fulfillmentOrderLineItems,omitempty"`
}

type FulfillmentOrderLineItemInput struct {
	ID       graphql.ID  `json:"id,omitempty"`
	Quantity graphql.Int `json:"quantity,omitempty"`
}

type FulfillmentTrackingInput struct {
	Company graphql.String `json:"company,omitempty"`
	Number  graphql.String `json:"number,omitempty"`
	URL     URL            `json:"url,omitempty"`
}

type mutationFulfillmentCreateV2 struct {
	FulfillmentCreateV2Result FulfillmentCreateV2Result `graphql:"fulfillmentCreateV2(fulfillment: $fulfillment)" json:"fulfillmentCreateV2"`
}

type FulfillmentCreateV2Result struct {
	UserErrors []UserErrors `json:"userErrors,omitempty"`
}

func (s *FulfillmentServiceOp) Create(fulfillment FulfillmentV2Input) error {
	m := mutationFulfillmentCreateV2{}

	vars := map[string]interface{}{
		"fulfillment": fulfillment,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("Mutation error: %s", err)
	}

	if len(m.FulfillmentCreateV2Result.UserErrors) > 0 {
		return fmt.Errorf("UserErrors: %+v", m.FulfillmentCreateV2Result.UserErrors)
	}

	return nil
}
