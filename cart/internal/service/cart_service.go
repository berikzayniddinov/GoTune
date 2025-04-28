package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gotune/cart/internal/repository"
	"gotune/cart/proto"
	"gotune/events"
)

type CartService struct {
	repo repository.CartRepository
	proto.UnimplementedCartServiceServer
	eventPublisher *events.EventPublisher
}

func NewCartService(repo repository.CartRepository, publisher *events.EventPublisher) *CartService {
	return &CartService{repo: repo, eventPublisher: publisher}
}

func (s *CartService) AddToCart(ctx context.Context, req *proto.AddToCartRequest) (*proto.AddToCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}
	instrumentID, err := primitive.ObjectIDFromHex(req.InstrumentId)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddToCart(ctx, userID, instrumentID, req.Quantity); err != nil {
		return nil, err
	}

	_ = s.eventPublisher.Publish("cart_updated", map[string]string{
		"user_id":       req.UserId,
		"instrument_id": req.InstrumentId,
	})

	return &proto.AddToCartResponse{Success: true}, nil
}

func (s *CartService) GetCart(ctx context.Context, req *proto.GetCartRequest) (*proto.GetCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
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

	return &proto.GetCartResponse{Items: protoItems}, nil
}

func (s *CartService) RemoveFromCart(ctx context.Context, req *proto.RemoveFromCartRequest) (*proto.RemoveFromCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}
	instrumentID, err := primitive.ObjectIDFromHex(req.InstrumentId)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RemoveFromCart(ctx, userID, instrumentID); err != nil {
		return nil, err
	}

	return &proto.RemoveFromCartResponse{Success: true}, nil
}

func (s *CartService) ClearCart(ctx context.Context, req *proto.ClearCartRequest) (*proto.ClearCartResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	return &proto.ClearCartResponse{Success: true}, nil
}
