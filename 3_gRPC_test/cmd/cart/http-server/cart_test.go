package httpserver

import (
	"MyMod/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartServicePostValidSku(t *testing.T) {
	t.Parallel()

	mockClient := &MockGRPcClient{}
	service := NewCartService(mockClient)
	ui := NewUI(service)

	testCase := struct {
		UserID           int
		SkuID            int
		Count            uint32
		ExpectedPostCode int
		ExpectedGetCode  int
		ExpectedName     string
		ExpectedPrice    uint32
	}{
		1,
		1076963,
		2,
		http.StatusOK,
		http.StatusOK,
		"Теория нравственных чувств | Смит Адам",
		3379,
	}

	reqPost := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/user/%d/cart/%d", testCase.UserID, testCase.SkuID),
		strings.NewReader(fmt.Sprintf(`{"count": %d}`, testCase.Count)),
	)
	reqPost.Header.Set("Content-Type", "application/json")
	reqPost.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))
	reqPost.SetPathValue("sku_id", fmt.Sprintf("%d", testCase.SkuID))

	recPost := httptest.NewRecorder()
	ui.PostSKU(recPost, reqPost)

	require.Equal(t, testCase.ExpectedPostCode, recPost.Code)

	reqGet := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/user/%d/cart", testCase.UserID),
		nil,
	)
	reqGet.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))

	recGet := httptest.NewRecorder()

	ui.GetCart(recGet, reqGet)

	require.Equal(t, testCase.ExpectedGetCode, recGet.Code)

	var cartResponse domain.CartResponse
	err := json.NewDecoder(recGet.Body).Decode(&cartResponse)
	require.NoError(t, err)

	require.NotEmpty(t, cartResponse.Items)

	item := cartResponse.Items[0]

	assert.Equal(t, testCase.ExpectedName, item.Name)
	assert.Equal(t, testCase.ExpectedPrice, item.Price)
	assert.Equal(t, testCase.Count, item.Count)
}

func TestCartServicePostInvalidSku(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		TestName     string
		UserID       int
		SkuID        int
		Count        uint32
		ExpectedCode int
	}{
		{
			"zero sku",
			1,
			0,
			1,
			http.StatusBadRequest,
		},
		{
			"negative sku",
			1,
			-1,
			1,
			http.StatusBadRequest,
		},
		{
			"not exist sku",
			1,
			77767,
			1,
			http.StatusPreconditionFailed,
		},
	} {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			mockClient := &MockGRPcClient{}
			service := NewCartService(mockClient)
			ui := NewUI(service)

			reqPost := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/user/%d/cart/%d", testCase.UserID, testCase.SkuID),
				strings.NewReader(fmt.Sprintf(`{"count": %d}`, testCase.Count)),
			)
			reqPost.Header.Set("Content-Type", "application/json")
			reqPost.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))
			reqPost.SetPathValue("sku_id", fmt.Sprintf("%d", testCase.SkuID))

			recPost := httptest.NewRecorder()
			ui.PostSKU(recPost, reqPost)

			require.Equal(t, testCase.ExpectedCode, recPost.Code)
		})
	}
}

func TestCartServicePostSumCount(t *testing.T) {
	t.Parallel()

	mockClient := &MockGRPcClient{}
	service := NewCartService(mockClient)
	ui := NewUI(service)

	testCase := struct {
		UserID           int
		SkuID            int
		Count            uint32
		ExpectedPostCode int
		ExpectedGetCode  int
		ExpectedName     string
		ExpectedPrice    uint32
	}{
		1,
		1076963,
		2,
		http.StatusOK,
		http.StatusOK,
		"Теория нравственных чувств | Смит Адам",
		3379,
	}
	sumCount := 0
	for i := 1; i <= int(testCase.Count); i++ {
		sumCount += i
		reqPost := httptest.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/user/%d/cart/%d", testCase.UserID, testCase.SkuID),
			strings.NewReader(fmt.Sprintf(`{"count": %d}`, i)),
		)
		reqPost.Header.Set("Content-Type", "application/json")
		reqPost.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))
		reqPost.SetPathValue("sku_id", fmt.Sprintf("%d", testCase.SkuID))

		recPost := httptest.NewRecorder()
		ui.PostSKU(recPost, reqPost)

		require.Equal(t, testCase.ExpectedPostCode, recPost.Code, "request count:", i)
	}

	reqGet := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/user/%d/cart", testCase.UserID),
		nil,
	)
	reqGet.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))

	recGet := httptest.NewRecorder()

	ui.GetCart(recGet, reqGet)

	require.Equal(t, testCase.ExpectedGetCode, recGet.Code)

	var cartResponse domain.CartResponse
	err := json.NewDecoder(recGet.Body).Decode(&cartResponse)
	require.NoError(t, err)

	require.NotEmpty(t, cartResponse.Items)

	item := cartResponse.Items[0]

	assert.Equal(t, testCase.ExpectedName, item.Name)
	assert.Equal(t, testCase.ExpectedPrice, item.Price)
	assert.EqualValues(t, sumCount, item.Count)
	assert.EqualValues(t, sumCount*int(testCase.ExpectedPrice), cartResponse.TotalPrice)
}

func TestCartServiceGetInvalidUser(t *testing.T) {
	for _, testCase := range []struct {
		TestName     string
		UserID       int
		ExpectedCode int
	}{
		{
			"negative user_id",
			-1,
			http.StatusBadRequest,
		},
		{
			"not exist user_id",
			1,
			http.StatusNotFound,
		},
	} {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			mockClient := &MockGRPcClient{}
			service := NewCartService(mockClient)
			ui := NewUI(service)

			reqGet := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/user/%d/cart", testCase.UserID),
				nil,
			)
			reqGet.SetPathValue("user_id", fmt.Sprintf("%d", testCase.UserID))

			recGet := httptest.NewRecorder()

			ui.GetCart(recGet, reqGet)

			require.Equal(t, testCase.ExpectedCode, recGet.Code)

		})
	}
}

type MockGRPcClient struct {
}

func (c *MockGRPcClient) GetProduct(ctx context.Context, skuId uint32) (domain.Product, error) {
	switch skuId {
	case 1076963:
		return domain.Product{
			Name:  "Теория нравственных чувств | Смит Адам",
			Price: 3379,
		}, nil
	case 0:
		return domain.Product{}, domain.ErrInvalidArgument
	default:
		return domain.Product{}, domain.ErrPreconditionFailed
	}
}
