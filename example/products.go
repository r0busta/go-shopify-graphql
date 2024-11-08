package main

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/v4/graph/model"
	"github.com/r0busta/go-shopify-graphql/v9"
	"gopkg.in/guregu/null.v4"
)

func listProducts(client *shopify.Client) {
	// Get products
	products, err := client.Product.List(context.Background(), "")
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, p := range products {
		fmt.Println(p.Title)
	}
}

func createProduct(client *shopify.Client) {
	status := model.ProductStatusDraft
	input := model.ProductCreateInput{
		Title:  model.NewString("go-shopify-graphql T-Shirt"),
		Handle: model.NewString("go-shopify-graphql-t-shirt"),
		Status: &status,
		ProductOptions: []model.OptionCreateInput{
			{
				Name: model.NewString("Color"),
				Values: []model.OptionValueCreateInput{
					{
						Name: model.NewString("Red"),
					},
					{
						Name: model.NewString("Blue"),
					},
				},
			},
			{
				Name: model.NewString("Size"),
				Values: []model.OptionValueCreateInput{
					{
						Name: model.NewString("Small"),
					},
					{
						Name: model.NewString("Medium"),
					},
					{
						Name: model.NewString("Large"),
					},
				},
			},
		},
	}

	media := []model.CreateMediaInput{
		{
			OriginalSource:   "https://picsum.photos/seed/c854bed7-d604-4b8f-a5d3-0c44bb01a534/600/300",
			MediaContentType: model.MediaContentTypeImage,
		},
		{
			OriginalSource:   "https://picsum.photos/seed/3474e3f2-7ffe-4349-89f9-0deca2a87986/600/300",
			MediaContentType: model.MediaContentTypeImage,
		},
	}

	id, err := client.Product.Create(context.Background(), input, media)
	if err != nil {
		panic(err)
	}

	product, err := client.Product.Get(context.Background(), *id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created product", product.Title, product.ID)

	variants := []model.ProductVariantsBulkInput{{
		InventoryItem: &model.InventoryItemInput{
			Sku: model.NewString("t-shirt-1-red-small"),
		},
		Price: model.NewNullString(null.StringFrom("10.00")),
		OptionValues: []model.VariantOptionValueInput{
			{
				OptionName: model.NewString("Color"),
				Name:       model.NewString("Red"),
			},
			{
				OptionName: model.NewString("Size"),
				Name:       model.NewString("Small"),
			},
		},
	}, {
		InventoryItem: &model.InventoryItemInput{
			Sku: model.NewString("t-shirt-1-blue-large"),
		},
		Price: model.NewNullString(null.StringFrom("10.00")),
		OptionValues: []model.VariantOptionValueInput{
			{
				OptionName: model.NewString("Color"),
				Name:       model.NewString("Blue"),
			},
			{
				OptionName: model.NewString("Size"),
				Name:       model.NewString("Large"),
			},
		},
	}}

	err = client.Product.VariantsBulkCreate(context.Background(), *id, variants, model.ProductVariantsBulkCreateStrategyRemoveStandaloneVariant)
	if err != nil {
		panic(err)
	}

	product, err = client.Product.Get(context.Background(), *id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Added", len(product.Variants.Edges), "variants")

	err = client.Product.Delete(context.Background(), model.ProductDeleteInput{
		ID: *id,
	})
	if err != nil {
		panic(err)
	}
}
