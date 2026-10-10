package domain

import "errors"

type Products map[uint32]*Product

type Product struct {
	Name  string
	Price uint32
}

var (
	ErrNotFound           error = errors.New("not found")
	ErrInvalidArgument    error = errors.New("wrong id")
	InternalErr           error = errors.New("server error")
	ErrPreconditionFailed error = errors.New("sku not exist")
)

type CartItem struct {
	Count uint32
	Name  string
	Price uint32
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}
type CartItemResponse struct {
	SKU   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}
