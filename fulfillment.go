package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/v2/graph/model"
)

//go:generate mockgen -destination=./mock/fulfillment_service.go -package=mock . FulfillmentService
type FulfillmentService interface {
	Create(input model.FulfillmentV2Input) error
}

type FulfillmentServiceOp struct {
	client *Client
}

var _ FulfillmentService = &FulfillmentServiceOp{}

type mutationFulfillmentCreateV2 struct {
	FulfillmentCreateV2Result struct {
		UserErrors []model.UserError `json:"userErrors,omitempty"`
	} `graphql:"fulfillmentCreateV2(fulfillment: $fulfillment)" json:"fulfillmentCreateV2"`
}

func (s *FulfillmentServiceOp) Create(fulfillment model.FulfillmentV2Input) error {
	m := mutationFulfillmentCreateV2{}

	vars := map[string]interface{}{
		"fulfillment": fulfillment,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("mutation: %s", err)
	}

	if len(m.FulfillmentCreateV2Result.UserErrors) > 0 {
		return fmt.Errorf("UserErrors: %+v", m.FulfillmentCreateV2Result.UserErrors)
	}

	return nil
}
