package main

import (
	"fmt"

	shopify "github.com/r0busta/go-shopify-graphql/v7"
)

func collections(client *shopify.Client) {
	// Get all collections
	collections, err := client.Collection.ListAll()
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, c := range collections {
		fmt.Println(c.Handle)
	}
}
