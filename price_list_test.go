package shopify

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceListServiceOp_GetAll(t *testing.T) {
	t.Run("should return error if it fails to retrieve price lists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBulkOperationService := NewMockBulkOperationService(ctrl)

		mockBulkOperationService.EXPECT().
			BulkQuery(priceListsGetAllBulkQuery, gomock.Any()).
			Return(errors.New("error")).
			Times(1)

		priceListService := PriceListServiceOp{service: mockBulkOperationService}

		_, err := priceListService.GetAll()
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

		mockBulkOperationService := NewMockBulkOperationService(ctrl)

		mockBulkOperationService.EXPECT().
			BulkQuery(priceListsGetAllBulkQuery, gomock.Any()).
			SetArg(1, expectedPriceLists).
			Return(nil).
			Times(1)

		priceListService := PriceListServiceOp{service: mockBulkOperationService}

		priceLists, err := priceListService.GetAll()
		require.NoError(t, err)

		assert.Equal(t, expectedPriceLists, priceLists)
	})
}
