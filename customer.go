package shopify

import "github.com/es-hs/go-shopify-graphql/graphql"

type Customer struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
	FirstName        graphql.String `json:"firstName,omitempty"`
	DisplayName      graphql.String `json:"displayName,omitempty"`
	Email            graphql.String `json:"email,omitempty"`
}
