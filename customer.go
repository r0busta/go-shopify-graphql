package shopify

import "github.com/shurcooL/graphql"

type Customer struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
}
