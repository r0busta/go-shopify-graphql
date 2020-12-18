package shopify

import "github.com/r0busta/graphql"

type InventoryItem struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
}
