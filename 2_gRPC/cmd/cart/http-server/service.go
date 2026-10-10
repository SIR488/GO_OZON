package httpserver

import (
	"MyMod/domain"
	"context"
	"fmt"
	"sort"
	"time"
)

type Client interface {
	GetProduct(ctx context.Context, skuId uint32) (domain.Product, error)
}

type CartService struct {
	productClient Client
	carts         map[uint32]map[uint32]*domain.CartItem
}

func NewCartService(client Client) *CartService {
	service := &CartService{productClient: client}
	service.carts = make(map[uint32]map[uint32]*domain.CartItem)
	return service
}

func (s *CartService) InsertSKU(userId, skuId, count uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product, err := s.productClient.GetProduct(ctx, skuId)
	if err != nil {
		return fmt.Errorf("sku check error: %w", err)
	}
	_, ok := s.carts[userId]
	if !ok {
		s.carts[userId] = make(map[uint32]*domain.CartItem)
	}
	_, ok = s.carts[userId][skuId]
	if ok {
		s.carts[userId][skuId].Count += count
	} else {
		s.carts[userId][skuId] = &domain.CartItem{Count: count, Name: product.Name, Price: product.Price}
	}
	return nil
}
func (s *CartService) GetCart(userId uint32) (domain.CartResponse, error) {
	cart, ok := s.carts[userId]
	if !ok {
		return domain.CartResponse{}, domain.ErrNotFound
	}

	items := make([]domain.CartItemResponse, 0, len(cart))
	var total uint32

	for sku, item := range cart {
		items = append(items, domain.CartItemResponse{
			SKU:   int64(sku),
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
		total += item.Price * item.Count
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})

	return domain.CartResponse{
		Items:      items,
		TotalPrice: total,
	}, nil
}
