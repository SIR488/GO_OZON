package httpserver

import (
	"MyMod/domain"
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

type Client interface {
	GetProduct(ctx context.Context, skuId uint32) (domain.Product, error)
}

type CartService struct {
	productClient Client
	carts         map[uint32]map[uint32]*domain.CartItem
	mu            sync.RWMutex
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
	s.mu.Lock()
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
	s.mu.Unlock()
	return nil
}
func (s *CartService) GetCart(userId uint32) (domain.CartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.carts[userId]
	if !ok {
		return domain.CartResponse{}, domain.ErrNotFound
	}

	items := make([]domain.CartItemResponse, len(cart))
	var total atomic.Uint32
	i := 0

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	wg, ctx := errgroup.WithContext(ctxTimeout)
	defer cancel()

	for skuCycle, item := range cart {
		if ctx.Err() != nil {
			break
		}
		idx := i
		i++
		sku, count := skuCycle, item.Count
		wg.Go(func() error {
			product, err := s.productClient.GetProduct(ctx, sku)
			if err != nil {
				return err
			}
			items[idx] = domain.CartItemResponse{
				SKU:   int64(sku),
				Name:  product.Name,
				Count: count,
				Price: product.Price,
			}
			total.Add(count * product.Price)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return domain.CartResponse{}, errors.Join(domain.ErrPreconditionFailed, err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})

	return domain.CartResponse{
		Items:      items,
		TotalPrice: total.Load(),
	}, nil
}
