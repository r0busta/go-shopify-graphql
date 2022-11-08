package shopify_test

import (
	"os"
	"testing"

	"github.com/r0busta/go-shopify-graphql-model/v3/graph/model"
	"github.com/r0busta/go-shopify-graphql/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkOperationEndToEnd(t *testing.T) {
	require.NotZero(t, os.Getenv("STORE_API_KEY"))
	require.NotZero(t, os.Getenv("STORE_PASSWORD"))
	require.NotZero(t, os.Getenv("STORE_NAME"))
	require.NotZero(t, os.Getenv("STORE_ACCESS_TOKEN"))

	tests := []struct {
		name   string
		client *shopify.Client
	}{{
		name:   "default client",
		client: shopify.NewDefaultClient(),
	}, {
		name:   "client with a token",
		client: shopify.NewClientWithToken(os.Getenv("STORE_ACCESS_TOKEN"), os.Getenv("STORE_NAME")),
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
										media{
											edges {
												node {
													... on MediaImage {
														id
														image {
															url
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`

			res := []*model.Product{}
			err := tt.client.BulkOperation.BulkQuery(q, &res)
			require.NoError(t, err)

			assert.Greater(t, len(res), 1)
			assert.NotZero(t, res[0].ID)

			assert.Greater(t, len(res[0].Variants.Edges), 1)
			assert.NotZero(t, res[0].Variants.Edges[0].Node.ID)

			assert.Equal(t, len(res[0].Variants.Edges[0].Node.Media.Edges), 1)
			assert.NotZero(t, res[0].Variants.Edges[0].Node.Media.Edges[0].Node.(*model.MediaImage).ID)
			assert.NotEmpty(t, res[0].Variants.Edges[0].Node.Media.Edges[0].Node.(*model.MediaImage).Image.URL)
		})
	}
}
