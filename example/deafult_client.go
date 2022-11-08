package main

import (
	shopify "github.com/r0busta/go-shopify-graphql/v7"
)

func defaultClient() *shopify.Client {
	return shopify.NewDefaultClient()
}
