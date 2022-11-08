package main

import (
	"os"

	shopify "github.com/r0busta/go-shopify-graphql/v6"
)

func clientWithToken() {
	// Create client
	client := shopify.NewClientWithToken(os.Getenv("STORE_ACCESS_TOKEN"), os.Getenv("STORE_NAME"))

	// Collections
	collections(client)

	// Products
	products(client)

	// Bulk operations
	bulk(client)
}
