package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gotune/cart/internal/repository"
	"gotune/cart/metrics"
	"gotune/cart/proto"
	"gotune/events"
)

const (
	cartCacheKeyPrefix = "cart:"
	cartCacheTTL       = 10 * time.Minute
)

type CartService struct {
	repo           repository.CartRepository
	cache          *redis.Client
	eventPublisher *events.EventPublisher
	proto.UnimplementedCartServiceServer
}

func NewCartService(repo repository.CartRepository, publisher *events.EventPublisher, cache *redis.Client) *CartService {
	return &CartService{
		repo:           repo,
		cache:          cache,
		eventPublisher: publisher,
	}
}

func (s *CartService) AddToCart(ctx context.Context, req *proto.AddToCartRequest) (*proto.AddToCartResponse, error) {
	metrics.CartCreateAttempts.Inc()
	timer := prometheus.NewTimer(metrics.CartCreateDuration)
	defer timer.ObserveDuration()

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	instrumentID, err := primitive.ObjectIDFromHex(req.InstrumentId)
	if err != nil {
		return nil, fmt.Errorf("invalid instrument ID: %w", err)
	}

	if err := s.repo.AddToCart(ctx, userID, instrumentID, req.Quantity); err != nil {
		return nil, err
	}

	s.invalidateCartCache(ctx, req.UserId)

	_ = s.eventPublisher.Publish("cart_updated", map[string]string{
		"user_id":       req.UserId,
		"instrument_id": req.InstrumentId,
	})

	metrics.CartCreatedTotal.Inc()

	return &proto.AddToCartResponse{Success: true}, nil
}

func (s *CartService) GetCart(ctx context.Context, req *proto.GetCartRequest) (*proto.GetCartResponse, error) {
	cacheKey := cartCacheKeyPrefix + req.UserId
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp proto.GetCartResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	items, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var protoItems []*proto.CartItem
	for _, item := range items {
		protoItems = append(protoItems, &proto.CartItem{
			InstrumentId: item.InstrumentID.Hex(),
			Quantity:     item.Quantity,
		})
	}

	resp := &proto.GetCartResponse{Items: protoItems}
	data, _ := json.Marshal(resp)
	_ = s.cache.Set(ctx, cacheKey, data, cartCacheTTL)

	return resp, nil
}

func (s *CartService) RemoveFromCart(ctx context.Context, req *proto.RemoveFromCartRequest) (*proto.RemoveFromCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	instrumentID, err := primitive.ObjectIDFromHex(req.InstrumentId)
	if err != nil {
		return nil, fmt.Errorf("invalid instrument ID: %w", err)
	}

	if err := s.repo.RemoveFromCart(ctx, userID, instrumentID); err != nil {
		return nil, err
	}

	s.invalidateCartCache(ctx, req.UserId)

	_ = s.eventPublisher.Publish("cart_updated", map[string]string{
		"user_id":       req.UserId,
		"instrument_id": req.InstrumentId,
	})

	return &proto.RemoveFromCartResponse{Success: true}, nil
}

func (s *CartService) ClearCart(ctx context.Context, req *proto.ClearCartRequest) (*proto.ClearCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if err := s.repo.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	s.invalidateCartCache(ctx, req.UserId)

	_ = s.eventPublisher.Publish("cart_cleared", map[string]string{
		"user_id": req.UserId,
	})

	return &proto.ClearCartResponse{Success: true}, nil
}

func (s *CartService) DeleteAllCartCache(ctx context.Context, req *proto.DeleteAllCartCacheRequest) (*proto.DeleteAllCartCacheResponse, error) {
	keys, err := s.cache.Keys(ctx, cartCacheKeyPrefix+"*").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей кэша: %w", err)
	}

	if len(keys) > 0 {
		if err := s.cache.Del(ctx, keys...).Err(); err != nil {
			return nil, fmt.Errorf("ошибка очистки кэша: %w", err)
		}
	}

	return &proto.DeleteAllCartCacheResponse{Success: true}, nil
}

func (s *CartService) invalidateCartCache(ctx context.Context, userID string) {
	s.cache.Del(ctx, cartCacheKeyPrefix+userID)
}
