package service

import (
	"context"
	"gotune/order/internal/entity"
	"gotune/order/internal/repository"
	"gotune/order/proto"
	usersproto "gotune/users/proto"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	repo       repository.OrderRepository
	userClient usersproto.UserServiceClient
	proto.UnimplementedOrderServiceServer
}

func NewOrderService(repo repository.OrderRepository, userClient usersproto.UserServiceClient) *OrderService {
	return &OrderService{repo: repo, userClient: userClient}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// Проверка существования пользователя
	_, err := s.userClient.GetUser(ctx, &usersproto.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	userID, _ := primitive.ObjectIDFromHex(req.UserId)

	var items []entity.OrderItem
	for _, i := range req.Items {
		instrumentID, _ := primitive.ObjectIDFromHex(i.InstrumentId)
		items = append(items, entity.OrderItem{
			InstrumentID: instrumentID,
			Quantity:     i.Quantity,
		})
	}

	order := &entity.Order{
		UserID:    userID,
		Items:     items,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &proto.CreateOrderResponse{
		OrderId: id.Hex(),
	}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, req *proto.GetOrdersRequest) (*proto.GetOrdersResponse, error) {
	// Проверка существования пользователя
	_, err := s.userClient.GetUser(ctx, &usersproto.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	userID, _ := primitive.ObjectIDFromHex(req.UserId)

	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var protoOrders []*proto.Order
	for _, o := range orders {
		var protoItems []*proto.OrderItem
		for _, i := range o.Items {
			protoItems = append(protoItems, &proto.OrderItem{
				InstrumentId: i.InstrumentID.Hex(),
				Quantity:     i.Quantity,
			})
		}
		protoOrders = append(protoOrders, &proto.Order{
			OrderId:   o.ID.Hex(),
			Items:     protoItems,
			CreatedAt: o.CreatedAt.Unix(),
		})
	}

	return &proto.GetOrdersResponse{
		Orders: protoOrders,
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	// Проверка существования пользователя
	_, err := s.userClient.GetUser(ctx, &usersproto.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	orderID, _ := primitive.ObjectIDFromHex(req.OrderId)
	userID, _ := primitive.ObjectIDFromHex(req.UserId)

	if err := s.repo.Delete(ctx, orderID, userID); err != nil {
		return nil, err
	}

	return &proto.DeleteOrderResponse{Success: true}, nil
}
