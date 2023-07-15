package main

import (
	shopify "github.com/r0busta/go-shopify-graphql/v8"
)

func defaultClient() *shopify.Client {
	return shopify.NewDefaultClient()
}
