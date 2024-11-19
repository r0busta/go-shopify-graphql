package shopify

import (
	"os"

	"github.com/r0busta/graphql"
	log "github.com/sirupsen/logrus"
	graphqlclient "github.com/xiatechs/go-shopify-graphql/v4/graphql"
)

const (
	shopifyAPIVersion = "2024-10"
)

type Client struct {
	gql *graphql.Client

	Product      ProductService
	Variant      VariantService
	Inventory    InventoryService
	Collection   CollectionService
	Order        OrderService
	Fulfillment  FulfillmentService
	Location     LocationService
	Metafield    MetafieldService
	PriceList    PriceListService
	StagedUpload StagedUploadService
	File         FileService

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

// NewClient returns a new client based on the configuration arguments received.
// Most of this client functionality is not enabled since it has no test coverage.
// Although it should, we cannot be sure that it works on the current api version since we haven't tested it.
// Enabling should require adding test coverage and proper testing.
// Calling any of the disabled (commented below) functionaly will likely result in a nil pointer de-reference.
func NewClient(apiKey string, password string, storeName string) *Client {
	c := &Client{gql: newShopifyGraphQLClient(apiKey, password, storeName)}

	c.BulkOperation = &BulkOperationServiceOp{client: c}

	// c.Product = &ProductServiceOp{client: c}
	// c.Variant = &VariantServiceOp{client: c}
	// c.Inventory = &InventoryServiceOp{client: c}
	// c.Collection = &CollectionServiceOp{client: c}
	// c.Order = &OrderServiceOp{client: c}
	// c.Fulfillment = &FulfillmentServiceOp{client: c}
	// c.Location = &LocationServiceOp{client: c}
	// c.Metafield = &MetafieldServiceOp{client: c}

	c.PriceList = &PriceListServiceOp{c.BulkOperation, c.gql}
	c.StagedUpload = &StagedUploadOp{mutationClient: c.gql}
	c.File = &FileOp{gqlClient: c.gql}

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
