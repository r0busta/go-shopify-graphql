# go-shopify-graphql

Simple go client for the Shopify's GraphQL Admin API.

The client is just a wrapper for the `github.com/shurcooL/graphql`, so check the usage docs at [github.com/shurcooL/graphql](github.com/shurcooL/graphql)

## Getting started

A Hello World example

```go
package main

import (
    "fmt"

    shopifygraphql "github.com/r0busta/go-shopify-graphql"
    "github.com/shurcooL/graphql"
)

func main(){
    // Create client
    opts := []shopifygraphql.Option{
        shopifygraphql.WithVersion("2019-10"),
        shopifygraphql.WithPrivateAppAuth(<YOUR_PRIVATE_APP_KEY>, <YOUR_PRIVATE_APP_PASSWORD>),
    }

    client := shopifygraphql.NewClient(<YOUR_STORE_NAME>, opts...)

    // Get first 100 products
    var q struct {
        Products struct {
            Edges []struct {
                Node struct {
                    ID    graphql.ID
                    Title graphql.String
                }
            }
        } `graphql:"products(first: $first)"`
    }
    vars := map[string]interface{}{
        "first": graphql.Int(100),
    }

    err := client.Query(context.Background(), &q, vars)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v", q)
}
```
