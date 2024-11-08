package main

import (
	"os"

	shopify "github.com/r0busta/go-shopify-graphql/v9"
	graphqlclient "github.com/r0busta/go-shopify-graphql/v9/graphql"
)

func defaultClient() *shopify.Client {
	if os.Getenv("STORE_API_KEY") == "" || os.Getenv("STORE_PASSWORD") == "" || os.Getenv("STORE_NAME") == "" {
		panic("Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	if os.Getenv("STORE_API_VERSION") != "" {
		apiKey := os.Getenv("STORE_API_KEY")
		accessToken := os.Getenv("STORE_PASSWORD")
		storeName := os.Getenv("STORE_NAME")
		opts := []graphqlclient.Option{
			graphqlclient.WithVersion(os.Getenv("STORE_API_VERSION")),
			graphqlclient.WithPrivateAppAuth(apiKey, accessToken),
		}

		gql := graphqlclient.NewClient(storeName, opts...)

		return shopify.NewClient(shopify.WithGraphQLClient(gql))
	}

	return shopify.NewDefaultClient()
}
