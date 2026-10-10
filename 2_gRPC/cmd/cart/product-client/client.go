package client

import (
	"MyMod/domain"
	"MyMod/pkg"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductClient struct {
	innerClient pkg.ProductClient
}

func NewProductClient(client pkg.ProductClient) *ProductClient {
	return &ProductClient{innerClient: client}
}

func (c *ProductClient) GetProduct(ctx context.Context, skuId uint32) (domain.Product, error) {
	request := pkg.GetProductRequest{Sku: skuId}
	response, err := c.innerClient.GetProduct(ctx, &request)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return domain.Product{}, domain.ErrPreconditionFailed
		case codes.InvalidArgument:
			return domain.Product{}, domain.ErrInvalidArgument
		default:
			return domain.Product{}, fmt.Errorf("%w: %w", domain.InternalErr, err)
		}
	}

	return domain.Product{
		Name:  response.Name,
		Price: uint32(response.Price),
	}, nil
}
