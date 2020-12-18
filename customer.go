package shopify

import "github.com/r0busta/graphql"

type Customer struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
}
