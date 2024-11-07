package shopify

import (
	"os"

	graphqlclient "github.com/r0busta/go-shopify-graphql/v8/graphql"
	"github.com/r0busta/graphql"
	log "github.com/sirupsen/logrus"
)

const (
	defaultShopifyAPIVersion = "2025-01"
)

type Client struct {
	gql graphql.GraphQL

	Product       ProductService
	Inventory     InventoryService
	Collection    CollectionService
	Order         OrderService
	Fulfillment   FulfillmentService
	Location      LocationService
	Metafield     MetafieldService
	BulkOperation BulkOperationService
}

type Option func(shopClient *Client)

func WithGraphQLClient(gql graphql.GraphQL) Option {
	return func(c *Client) {
		c.gql = gql
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.gql == nil {
		log.Fatalln("GraphQL client not set")
	}

	c.Product = &ProductServiceOp{client: c}
	c.Inventory = &InventoryServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Fulfillment = &FulfillmentServiceOp{client: c}
	c.Location = &LocationServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.BulkOperation = &BulkOperationServiceOp{client: c}

	return c
}

func NewDefaultClient() *Client {
	apiKey := os.Getenv("STORE_API_KEY")
	accessToken := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API Key and/or Password (aka access token) and/or store name not set")
	}

	gql := newShopifyGraphQLClientWithBasicAuth(apiKey, accessToken, storeName)

	return NewClient(WithGraphQLClient(gql))
}

func NewClientWithToken(accessToken string, storeName string, opts ...Option) *Client {
	if accessToken == "" || storeName == "" {
		log.Fatalln("Shopify Admin API access token and/or store name not set")
	}

	gql := newShopifyGraphQLClientWithToken(accessToken, storeName)

	return NewClient(WithGraphQLClient(gql))
}

func newShopifyGraphQLClientWithBasicAuth(apiKey string, accessToken string, storeName string) *graphql.Client {
	opts := []graphqlclient.Option{
		graphqlclient.WithVersion(defaultShopifyAPIVersion),
		graphqlclient.WithPrivateAppAuth(apiKey, accessToken),
	}

	return graphqlclient.NewClient(storeName, opts...)
}

func newShopifyGraphQLClientWithToken(accessToken string, storeName string) *graphql.Client {
	opts := []graphqlclient.Option{
		graphqlclient.WithVersion(defaultShopifyAPIVersion),
		graphqlclient.WithToken(accessToken),
	}

	return graphqlclient.NewClient(storeName, opts...)
}

func (c *Client) GraphQLClient() graphql.GraphQL {
	return c.gql
}
