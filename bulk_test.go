package shopify_test

import (
	"os"
	"testing"

	"github.com/r0busta/go-shopify-graphql-model/v3/graph/model"
	"github.com/r0busta/go-shopify-graphql/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkOperationEndToEnd(t *testing.T) {
	require.NotZero(t, os.Getenv("STORE_API_KEY"))
	require.NotZero(t, os.Getenv("STORE_PASSWORD"))
	require.NotZero(t, os.Getenv("STORE_NAME"))

	client := shopify.NewDefaultClient()

	q := `
	{
		products{
			edges {
				node {
					id
					variants {
						edges {
							node {
								id
							}
						}
					}
				}
			}
		}
	}`

	res := []*model.Product{}
	err := client.BulkOperation.BulkQuery(q, &res)
	require.NoError(t, err)

	assert.Greater(t, len(res), 1)
	assert.NotZero(t, res[0].ID)

	assert.Greater(t, len(res[0].Variants.Edges), 1)
	assert.NotZero(t, res[0].Variants.Edges[0].Node.ID)
}
