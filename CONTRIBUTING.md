# Contributing

Thanks for helping out. This library started as a personal tool, so the
process is lightweight, but a few things make pull requests much easier to
merge.

## Before opening a pull request

- Run `gofmt`, `go vet ./...` and `go test ./...`. CI runs the same checks.
- Keep the module path as `github.com/r0busta/go-shopify-graphql/v10`. If you
  develop in a fork, use a `replace` directive locally instead of renaming the
  module in `go.mod`.
- Keep pull requests focused. One feature or fix per PR is ideal.
- If you add a service method, add the matching method to the interface and
  regenerate the mock in `mock/` (`go generate ./...`).

## Versioning

This client depends on
[go-shopify-graphql-model](https://github.com/r0busta/go-shopify-graphql-model),
which is generated from a specific Shopify Admin API schema version. Moving to
a new schema version is a breaking change in the model (new major version),
which in turn means a new major version of this client. Field additions to
existing queries that already exist in the current model are non-breaking and
can go into a minor or patch release.

## Testing against a store

The bulk operation end-to-end test needs a real store. It is skipped unless
`STORE_NAME`, `STORE_ACCESS_TOKEN`, `STORE_API_KEY` and `STORE_PASSWORD` are
set. See `.env.example`. Never commit `.env`.
