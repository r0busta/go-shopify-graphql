package shopify

import (
	graphqlclient "github.com/r0busta/go-shopify-graphql/v6/graphql"
	"github.com/r0busta/graphql"
	log "github.com/sirupsen/logrus"
)

const (
	defaultShopifyAPIVersion = "2022-04"
)

type Client struct {
	gql graphql.GraphQL

	Product       ProductService
	Variant       VariantService
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

func NewDefaultClient(apiKey, password, storeName string, opts ...Option) *Client {
	if apiKey == "" || password == "" || storeName == "" {
		log.Fatalln("Shopify app API Key and/or Password and/or Store Name not set")
	}

	return NewClient(apiKey, password, storeName, opts...)
}

func NewClient(apiKey string, password string, storeName string, opts ...Option) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.gql == nil {
		c.gql = newShopifyGraphQLClient(apiKey, password, storeName)
	}

	c.Product = &ProductServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Inventory = &InventoryServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Fulfillment = &FulfillmentServiceOp{client: c}
	c.Location = &LocationServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.BulkOperation = &BulkOperationServiceOp{client: c}

	return c
}

func newShopifyGraphQLClient(apiKey string, password string, storeName string) *graphql.Client {
	opts := []graphqlclient.Option{
		graphqlclient.WithVersion(defaultShopifyAPIVersion),
		graphqlclient.WithPrivateAppAuth(apiKey, password),
	}
	return graphqlclient.NewClient(storeName, opts...)
}

func (c *Client) GraphQLClient() graphql.GraphQL {
	return c.gql
}
