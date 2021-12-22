package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
)

//go:generate mockgen -destination=./mock/variant_service.go -package=mock . VariantService
type VariantService interface {
	Update(variant model.ProductVariantInput) error
}

type VariantServiceOp struct {
	client *Client
}

type mutationProductVariantUpdate struct {
	ProductVariantUpdateResult model.ProductVariantUpdatePayload `graphql:"productVariantUpdate(input: $input)" json:"productVariantUpdate"`
}

func (s *VariantServiceOp) Update(variant model.ProductVariantInput) error {
	m := mutationProductVariantUpdate{}

	vars := map[string]interface{}{
		"input": variant,
	}
	err := s.client.gql.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	if len(m.ProductVariantUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductVariantUpdateResult.UserErrors)
	}

	return nil
}
