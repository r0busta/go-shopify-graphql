package shopify

import (
	"context"

	"github.com/r0busta/graphql"
)

type LocationService interface {
	Get(id graphql.ID) (*Location, error)
}

type LocationServiceOp struct {
	client *Client
}

type Location struct {
	ID   graphql.ID     `json:"id,omitempty"`
	Name graphql.String `json:"name,omitempty"`
}

func (s *LocationServiceOp) Get(id graphql.ID) (*Location, error) {
	q := `query location($id: ID!) {
		location(id: $id){
			id
			name
		}
	}`

	vars := map[string]interface{}{
		"id": id,
	}

	out := struct {
		Location *Location `json:"location"`
	}{}
	err := s.client.gql.QueryString(context.Background(), q, vars, &out)
	if err != nil {
		return nil, err
	}

	return out.Location, nil
}
