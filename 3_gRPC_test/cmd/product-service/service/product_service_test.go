package service

import (
	"MyMod/pkg"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetProductGRPC(t *testing.T) {
	t.Parallel()

	service, err := NewProductService()
	require.NoError(t, err)

	ctx := context.TODO()

	for _, testCase := range []struct {
		TestName      string
		SkuID         uint32
		ExpectedName  string
		ExpectedPrice int32
		ExpectedCode  codes.Code
	}{
		{
			"valid sku",
			1076963,
			"Теория нравственных чувств | Смит Адам",
			3379,
			codes.OK,
		},
		{
			"not exist sku",
			77767,
			"",
			333,
			codes.NotFound,
		},
		{
			"zero value sku",
			0,
			"",
			333,
			codes.InvalidArgument,
		},
	} {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			request := pkg.GetProductRequest{Sku: testCase.SkuID}
			response, err := service.GetProduct(ctx, &request)
			if err != nil {
				require.Equal(t, testCase.ExpectedCode, status.Code(err))
				return
			}
			assert.Equal(t, testCase.ExpectedName, response.Name)
			assert.Equal(t, testCase.ExpectedPrice, response.Price)
		})

	}
}
