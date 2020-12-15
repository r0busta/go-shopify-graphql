package shopify

import (
	"log"
	"os"

	graphqlclient "github.com/r0busta/go-shopify-graphql/graphql"
	"github.com/shurcooL/graphql"
)

const (
	shopifyAPIVersion = "2021-01"
)

type Client struct {
	gql *graphql.Client

	Product       ProductService
	Variant       VariantService
	Collection    CollectionService
	Order         OrderService
	Metafield     MetafieldService
	BulkOperation BulkOperationService
}

type UserErrors struct {
	Field   []graphql.String
	Message graphql.String
}

func NewDefaultClient() (shopClient *Client) {
	apiKey := os.Getenv("STORE_API_KEY")
	password := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || password == "" || storeName == "" {
		log.Panicln("Shopify app API Key and/or Password and/or Store Name not set")
	}

	shopClient = NewClient(apiKey, password, storeName)

	return
}

func NewClient(apiKey string, password string, storeName string) *Client {
	c := &Client{gql: newShopifyGraphQLClient(apiKey, password, storeName)}

	c.Product = &ProductServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.BulkOperation = &BulkOperationServiceOp{client: c}

	return c
}

func newShopifyGraphQLClient(apiKey string, password string, storeName string) *graphql.Client {
	opts := []graphqlclient.Option{
		graphqlclient.WithVersion(shopifyAPIVersion),
		graphqlclient.WithPrivateAppAuth(apiKey, password),
	}
	return graphqlclient.NewClient(storeName, opts...)
}

func (c *Client) GraphQLClient() *graphql.Client {
	return c.gql
}
