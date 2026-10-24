package test

import (
	"MyMod/pkg"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var url string = "http://localhost:8080"
var client *http.Client = &http.Client{
	Timeout: 5 * time.Second,
}

type cartItemResponse struct {
	SKU   uint32 `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}

type cartResponse struct {
	Items      []cartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}

func generatePostRequest(userId, skuId int, count uint32) (*http.Request, error) {
	postURL := fmt.Sprintf("%s/user/%d/cart/%d", url, userId, skuId)
	postBody := strings.NewReader(fmt.Sprintf(`{"count": %d}`, count))

	reqPost, err := http.NewRequest(http.MethodPost, postURL, postBody)
	if err != nil {
		return nil, err
	}
	reqPost.Header.Set("Content-Type", "application/json")
	return reqPost, nil

}
func generateGetRequest(userId int) (*http.Request, error) {
	getURL := fmt.Sprintf("%s/user/%d/cart", url, userId)
	return http.NewRequest(http.MethodGet, getURL, nil)
}

func newGRPCTestClient(t *testing.T) pkg.ProductClient {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		conn.Close()
	})
	return pkg.NewProductClient(conn)
}
