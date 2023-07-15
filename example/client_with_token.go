package main

import (
	"os"

	shopify "github.com/r0busta/go-shopify-graphql/v8"
)

func clientWithToken() *shopify.Client {
	return shopify.NewClientWithToken(os.Getenv("STORE_ACCESS_TOKEN"), os.Getenv("STORE_NAME"))
}
