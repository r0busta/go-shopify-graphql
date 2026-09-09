# go-shopify-graphql

Go client for the Shopify GraphQL Admin API. Module path
`github.com/r0busta/go-shopify-graphql/v10`.

## Layout

- `client.go`: `Client` and constructors. `NewClientWithToken` is the modern
  path; `NewDefaultClient` is legacy basic auth for private apps.
- `graphql/`: thin wrapper that builds the Shopify endpoint URL and auth
  transport on top of `github.com/r0busta/graphql`.
- One file per Shopify resource (`product.go`, `order.go`, ...). Each defines
  a `XxxService` interface, an `XxxServiceOp` implementation, and the query
  strings it uses.
- `bulk.go`: bulk operations (submit, poll, download JSONL, decode).
- `mock/`: gomock mocks of the service interfaces, generated with
  `go generate ./...`. Regenerate after changing an interface.
- `example/`: runnable examples, not part of the library API.

## Dependencies to be aware of

- Types come from `github.com/r0busta/go-shopify-graphql-model/v5`. The
  default API version in `client.go` must match the schema that model version
  was generated from. Changing either is a coordinated major-version bump.
- `github.com/r0busta/graphql` is a fork of shurcooL/graphql maintained in
  the same GitHub account.

## Working here

- `go build ./... && go vet ./... && go test ./...` must pass cold. The bulk
  test skips itself without `STORE_*` env vars.
- `.env` holds live store credentials and is gitignored. Never stage it.
- Do not rename the module path or bump the `go` directive casually; consumers
  pin against both.
