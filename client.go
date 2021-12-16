package shopify

import (
	"os"

	"github.com/r0busta/graphql"
	log "github.com/sirupsen/logrus"
	graphqlclient "github.com/xiatechs/go-shopify-graphql/v4/graphql"
)

const (
	// TODO:
	//  I've had to bump this to 2021-04 (from 2021-01 which is what upstream supports) to support price lists
	//  I wonder if we should disabled all of the remaining functionality (besides price lists) that we don't know if it works with the new API version?
	//  There's also pretty much no coverage for all of the upstream functionality...
	shopifyAPIVersion = "2021-04"
)

type Client struct {
	gql *graphql.Client

	Product     ProductService
	Variant     VariantService
	Inventory   InventoryService
	Collection  CollectionService
	Order       OrderService
	Fulfillment FulfillmentService
	Location    LocationService
	Metafield   MetafieldService
	PriceList   PriceListService

	BulkOperation BulkOperationService
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

	bulkOperationService := &BulkOperationServiceOp{client: c}
	c.BulkOperation = bulkOperationService

	c.Product = &ProductServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Inventory = &InventoryServiceOp{client: c}
	c.Collection = &CollectionServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	c.Fulfillment = &FulfillmentServiceOp{client: c}
	c.Location = &LocationServiceOp{client: c}
	c.Metafield = &MetafieldServiceOp{client: c}
	c.PriceList = &PriceListServiceOp{bulkOperationService}

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
