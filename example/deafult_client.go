package main

import (
	shopify "github.com/r0busta/go-shopify-graphql/v6"
)

func defaultClient() *shopify.Client {
	return shopify.NewDefaultClient()
}
