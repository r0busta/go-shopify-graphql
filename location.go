package shopify

import (
	"context"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	"github.com/r0busta/graphql"
)

type LocationService interface {
	Get(id graphql.ID) (*model.Location, error)
}

type LocationServiceOp struct {
	client *Client
}

func (s *LocationServiceOp) Get(id graphql.ID) (*model.Location, error) {
	q := `query location($id: ID!) {
		location(id: $id){
			id
			name
		}
	}`

	vars := map[string]interface{}{
		"id": id,
	}

	var out struct {
		*model.Location `json:"location"`
	}
	err := s.client.gql.QueryString(context.Background(), q, vars, &out)
	if err != nil {
		return nil, err
	}

	return out.Location, nil
}
