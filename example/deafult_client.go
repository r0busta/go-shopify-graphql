package main

import (
	shopify "github.com/r0busta/go-shopify-graphql/v9"
)

func defaultClient() *shopify.Client {
	return shopify.NewDefaultClient()
}
