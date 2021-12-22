package shopify

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/r0busta/go-shopify-graphql-model/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestPriceListServiceOp_GetAll(t *testing.T) {
	t.Run("should return error if it fails to retrieve price lists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPriceListBulkQueryClient := NewMockPriceListBulkQueryClient(ctrl)

		mockPriceListBulkQueryClient.EXPECT().
			BulkQuery(gomock.Any(), gomock.Any()).
			Return(errors.New("error")).
			Times(1)

		priceListService := PriceListServiceOp{bulkQueryClient: mockPriceListBulkQueryClient}

		_, err := priceListService.GetPriceLists(context.Background())
		assert.Error(t, err)
	})

	t.Run("should return price lists successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedPriceLists := []PriceList{
			{
				ID:       "id",
				Name:     "name",
				Currency: "currency",
			},
		}

		mockPriceListBulkQueryClient := NewMockPriceListBulkQueryClient(ctrl)

		mockPriceListBulkQueryClient.EXPECT().
			BulkQuery(priceListsGetAllBulkQuery, gomock.Any()).
			SetArg(1, expectedPriceLists).
			Return(nil).
			Times(1)

		priceListService := PriceListServiceOp{bulkQueryClient: mockPriceListBulkQueryClient}

		priceLists, err := priceListService.GetPriceLists(context.Background())
		require.NoError(t, err)

		assert.Equal(t, expectedPriceLists, priceLists)
	})
}

func TestPriceListServiceOp_AddFixedPrice(t *testing.T) {
	var (
		ctx = context.Background()
	)

	t.Run("should return error if it fails to execute mutation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPriceListMutationClient := NewMockPriceListMutationClient(ctrl)

		mockPriceListMutationClient.EXPECT().
			Mutate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("error")).
			Times(1)

		priceListService := PriceListServiceOp{mutationClient: mockPriceListMutationClient}

		err := priceListService.AddFixedPricesToPriceList(ctx, "", nil)
		assert.Error(t, err)
	})

	t.Run("should add fixed price successfuly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		priceListID := "gid://shopify/PriceList/id"
		prices := []PriceListPriceInput{
			{
				CompareAtPrice: model.MoneyInput{
					Amount: null.String{
						NullString: sql.NullString{
							String: "10.2",
							Valid:  true,
						},
					},
					CurrencyCode: "EUR",
				},
				Price: model.MoneyInput{
					Amount: null.String{
						NullString: sql.NullString{
							String: "10.2",
							Valid:  true,
						},
					},
					CurrencyCode: "EUR",
				},
				VariantID: "gid://shopify/ProductVariant/id",
			},
		}

		mockPriceListMutationClient := NewMockPriceListMutationClient(ctrl)

		mockPriceListMutationClient.EXPECT().
			Mutate(
				ctx,
				&mutationPriceListFixedPricesAdd{},
				map[string]interface{}{
					"priceListId": priceListID,
					"prices":      prices,
				},
			).
			Return(nil).
			Times(1)

		priceListService := PriceListServiceOp{mutationClient: mockPriceListMutationClient}

		err := priceListService.AddFixedPricesToPriceList(ctx, priceListID, prices)
		assert.NoError(t, err)
	})
}
