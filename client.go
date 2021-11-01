package shopify

import (
	"os"

	graphqlclient "github.com/es-hs/go-shopify-graphql/graph"
	"github.com/es-hs/go-shopify-graphql/graphql"

	log "github.com/sirupsen/logrus"
)

const (
	shopifyAPIVersion = "2021-01"
)

type Client struct {
	gql *graphql.Client

	Product       ProductService
	Variant       VariantService
	Inventory     InventoryService
	Collection    CollectionService
	Order         OrderService
	Fulfillment   FulfillmentService
	Location      LocationService
	Metafield     MetafieldService
	BulkOperation BulkOperationService
	Webhook       WebhookService
}

type ListOptions struct {
	Query   string
	First   int
	Last    int
	After   string
	Before  string
	Reverse bool
}

func NewDefaultClient() (shopClient *Client) {
	apiKey := os.Getenv("STORE_API_KEY")
	password := os.Getenv("STORE_PASSWORD")
	storeName := os.Getenv("STORE_NAME")
	if apiKey == "" || password == "" || storeName == "" {
		log.Fatalln("Shopify app API Key and/or Password and/or Store Name not set")
	}

	shopClient = NewClient(apiKey, password, storeName)

	return
}

func NewClient(apiKey string, password string, storeName string) *Client {
	c := &Client{gql: newShopifyGraphQLClient(apiKey, password, storeName)}

	c.Product = &ProductServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Inventory = &InventoryServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Fulfillment = &FulfillmentServiceOp{client: c}
	c.Location = &LocationServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.BulkOperation = &BulkOperationServiceOp{client: c}
	c.Webhook = &WebhookServiceOp{client: c}

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

func NewClientWithToken(apiKey string, storeName string) *Client {
	c := &Client{gql: newShopifyGraphQLClientWithToken(apiKey, storeName)}

	c.Product = &ProductServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	// c.Inventory = &InventoryServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	// c.Order = &OrderServiceOp{client: c}
	// c.Fulfillment = &FulfillmentServiceOp{client: c}
	// c.Location = &LocationServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.BulkOperation = &BulkOperationServiceOp{client: c}
	c.Webhook = &WebhookServiceOp{client: c}

	return c
}

func newShopifyGraphQLClientWithToken(token string, storeName string) *graphql.Client {
	opts := []graphqlclient.Option{
		graphqlclient.WithVersion(shopifyAPIVersion),
		graphqlclient.WithToken(token),
	}
	// todo no more fixed storeName
	return graphqlclient.NewClient(storeName, opts...)
}
