package shopify

import "github.com/shurcooL/graphql"

type InventoryItem struct {
	ID               graphql.ID     `json:"id,omitempty"`
	LegacyResourceID graphql.String `json:"legacyResourceId,omitempty"`
}
