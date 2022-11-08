package main

import (
	"fmt"

	shopify "github.com/r0busta/go-shopify-graphql/v7"
)

func products(client *shopify.Client) {
	// Get products
	products, err := client.Product.List("")
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, p := range products {
		fmt.Println(p.Title)
	}
}
