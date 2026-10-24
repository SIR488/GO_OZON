package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainUserCase(t *testing.T) {
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
		1,
		http.StatusOK,
		http.StatusOK,
		"Теория нравственных чувств | Смит Адам",
		3379,
	}

	reqPost, err := generatePostRequest(testCase.UserID, testCase.SkuID, testCase.Count)
	require.NoError(t, err)
	respPost, err := client.Do(reqPost)
	require.NoError(t, err)
	defer respPost.Body.Close()

	require.Equal(t, testCase.ExpectedPostCode, respPost.StatusCode)

	reqGet, err := generateGetRequest(testCase.UserID)
	require.NoError(t, err)

	respGet, err := client.Do(reqGet)
	require.NoError(t, err)
	defer respGet.Body.Close()

	require.Equal(t, testCase.ExpectedGetCode, respGet.StatusCode)

	var cartResponse cartResponse
	err = json.NewDecoder(respGet.Body).Decode(&cartResponse)

	require.NoError(t, err)
	require.NotEmpty(t, cartResponse.Items)

	item := cartResponse.Items[0]

	assert.Equal(t, testCase.ExpectedName, item.Name)
	assert.Equal(t, testCase.ExpectedPrice, item.Price)
	assert.Equal(t, testCase.Count, item.Count)
}

func TestManyDifferentPosts(t *testing.T) {
	testCase := struct {
		UserID       int
		ExpectedCode int
	}{
		2,
		http.StatusOK,
	}

	postRequests := []struct {
		SkuID         int
		Count         uint32
		ExpectedName  string
		ExpectedPrice uint32
	}{
		{
			1076963,
			1,
			"Теория нравственных чувств | Смит Адам",
			3379,
		},
		{
			1148162,
			3,
			"Кулинар Гуров",
			2931,
		},
		{
			135147479,
			1,
			"Онлайн.ru. Сидишь в Интернете? Прочти! Тут о тебе! | Игошев Андрей",
			2981,
		},
	}
	var expectedTotalSum uint32 = 0
	for _, postRequest := range postRequests {

		reqPost, err := generatePostRequest(testCase.UserID, postRequest.SkuID, postRequest.Count)
		require.NoError(t, err)

		respPost, err := client.Do(reqPost)
		require.NoError(t, err)

		require.Equal(t, testCase.ExpectedCode, respPost.StatusCode)
		respPost.Body.Close()
		expectedTotalSum += postRequest.Count * postRequest.ExpectedPrice
	}

	reqGet, err := generateGetRequest(testCase.UserID)
	require.NoError(t, err)

	respGet, err := client.Do(reqGet)
	require.NoError(t, err)
	defer respGet.Body.Close()

	require.Equal(t, testCase.ExpectedCode, respGet.StatusCode)

	var cartResponse cartResponse
	err = json.NewDecoder(respGet.Body).Decode(&cartResponse)

	require.NoError(t, err)
	require.NotEmpty(t, cartResponse.Items)

	for i, item := range cartResponse.Items {
		assert.Equal(t, postRequests[i].ExpectedName, item.Name)
		assert.Equal(t, postRequests[i].ExpectedPrice, item.Price)
		assert.Equal(t, postRequests[i].Count, item.Count)
	}
	assert.Equal(t, expectedTotalSum, cartResponse.TotalPrice)
}

func TestGetUserWithEmptyCart(t *testing.T) {
	testCase := struct {
		UserID          int
		ExpectedGetCode int
	}{
		3,
		http.StatusNotFound,
	}

	reqGet, err := generateGetRequest(testCase.UserID)
	require.NoError(t, err)

	respGet, err := client.Do(reqGet)
	require.NoError(t, err)
	defer respGet.Body.Close()

	require.Equal(t, testCase.ExpectedGetCode, respGet.StatusCode)
}

func TestManyWrongPosts(t *testing.T) {
	userId := 4

	postRequests := []struct {
		SkuID        int
		Count        uint32
		ExpectedCode int
	}{
		{
			66676,
			1,
			http.StatusPreconditionFailed,
		},
		{
			-1148162,
			3,
			http.StatusBadRequest,
		},
		{
			0,
			1,
			http.StatusBadRequest,
		},
	}
	for _, postRequest := range postRequests {

		reqPost, err := generatePostRequest(userId, postRequest.SkuID, postRequest.Count)
		require.NoError(t, err)

		respPost, err := client.Do(reqPost)
		require.NoError(t, err)

		require.Equal(t, postRequest.ExpectedCode, respPost.StatusCode)
		respPost.Body.Close()
	}
}
