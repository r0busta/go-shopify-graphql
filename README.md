# go-shopify-graphql

A simple Go client for the Shopify GraphQL Admin API.

Objects and inputs come from
[go-shopify-graphql-model](https://github.com/r0busta/go-shopify-graphql-model),
which is generated from the Shopify Admin API schema. The client currently
targets API version `2026-07` by default. When Shopify retires a version, it
serves requests with the oldest supported version instead and reports the
version actually used in the `X-Shopify-API-Version` response header.

## Getting started

### 1. Setup

Create an Admin API access token for your store (Settings → Apps and sales
channels → Develop apps) and export it:

```bash
export STORE_NAME=<store_name>
export STORE_ACCESS_TOKEN=<access_token>
```

### 2. Program

```go
package main

import (
	"context"
	"fmt"
	"os"

	shopify "github.com/r0busta/go-shopify-graphql/v10"
)

func main() {
	client := shopify.NewClientWithToken(os.Getenv("STORE_ACCESS_TOKEN"), os.Getenv("STORE_NAME"))

	collections, err := client.Collection.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	for _, c := range collections {
		fmt.Println(c.Handle)
	}
}
```

### 3. Run

```bash
go run .
```

## Choosing the API version

To target a specific API version, build the underlying GraphQL client yourself:

```go
import (
	shopify "github.com/r0busta/go-shopify-graphql/v10"
	graphqlclient "github.com/r0busta/go-shopify-graphql/v10/graphql"
)

gql := graphqlclient.NewClient(storeName,
	graphqlclient.WithToken(accessToken),
	graphqlclient.WithVersion("2026-07"),
)
client := shopify.NewClient(shopify.WithGraphQLClient(gql))
```

Keep the version in step with the `go-shopify-graphql-model` release you
depend on, since the models are generated from a specific schema version.

## Legacy private app credentials

`NewDefaultClient` reads `STORE_API_KEY`, `STORE_PASSWORD` and `STORE_NAME`
and authenticates with HTTP basic auth, which is how private apps used to
work. New integrations should use `NewClientWithToken`.

## Raw queries and bulk operations

`client.GraphQLClient()` exposes the underlying GraphQL client for queries the
service methods do not cover. `client.BulkOperation.BulkQuery` runs a query as
a Shopify bulk operation and decodes the JSONL result. See the `example`
directory for both.

## Testing

```bash
go test ./...
```

The bulk operation test talks to a real store and is skipped unless
`STORE_NAME`, `STORE_ACCESS_TOKEN`, `STORE_API_KEY` and `STORE_PASSWORD` are
set. Copy `.env.example` to `.env` if you want to run it locally.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
