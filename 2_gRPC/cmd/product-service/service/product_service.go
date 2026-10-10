package service

import (
	"MyMod/domain"
	"MyMod/pkg"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pkg.ProductServer
	products domain.Products
}

func NewProductService() (*ProductService, error) {
	service := &ProductService{}
	products, err := domain.GetLocalData()
	if err != nil {
		return nil, fmt.Errorf("load products: %w", err)
	}
	service.products = products
	return service, nil
}

func (s *ProductService) GetProduct(ctx context.Context, in *pkg.GetProductRequest) (*pkg.GetProductResponse, error) {
	if in.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong id")
	}

	product, ok := s.products[in.Sku]

	if !ok {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return &pkg.GetProductResponse{Name: product.Name, Price: int32(product.Price)}, nil
}
