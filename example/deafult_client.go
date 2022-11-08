package main

import (
	shopify "github.com/r0busta/go-shopify-graphql/v6"
)

func defaultClient() {
	// Create client
	client := shopify.NewDefaultClient()

	// Collections
	collections(client)

	// Products
	products(client)

	// Bulk operations
	bulk(client)
}
