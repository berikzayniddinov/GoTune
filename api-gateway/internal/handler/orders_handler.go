package handler

import (
	"context"
	"encoding/json"
	"gotune/order/proto"
	usersproto "gotune/users/proto"
	"net/http"
)

type OrderHandler struct {
	OrderClient proto.OrderServiceClient
	UserClient  usersproto.UserServiceClient
}

func NewOrderHandler(orderClient proto.OrderServiceClient, userClient usersproto.UserServiceClient) *OrderHandler {
	return &OrderHandler{
		OrderClient: orderClient,
		UserClient:  userClient,
	}
}

type CreateOrderRequest struct {
	UserID string `json:"user_id"`
	Items  []struct {
		InstrumentID string `json:"instrument_id"`
		Quantity     int32  `json:"quantity"`
	} `json:"items"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка пользователя
	_, err := h.UserClient.GetUser(context.Background(), &usersproto.GetUserRequest{UserId: req.UserID})
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	var items []*proto.OrderItem
	for _, item := range req.Items {
		items = append(items, &proto.OrderItem{
			InstrumentId: item.InstrumentID,
			Quantity:     item.Quantity,
		})
	}

	_, err = h.OrderClient.CreateOrder(context.Background(), &proto.CreateOrderRequest{
		UserId: req.UserID,
		Items:  items,
	})
	if err != nil {
		http.Error(w, "Ошибка создания заказа", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Отсутствует user_id", http.StatusBadRequest)
		return
	}

	// Проверка пользователя
	_, err := h.UserClient.GetUser(context.Background(), &usersproto.GetUserRequest{UserId: userID})
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	resp, err := h.OrderClient.GetOrders(context.Background(), &proto.GetOrdersRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Ошибка получения заказов", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Orders)
}

type DeleteOrderRequest struct {
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	var req DeleteOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка пользователя
	_, err := h.UserClient.GetUser(context.Background(), &usersproto.GetUserRequest{UserId: req.UserID})
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	_, err = h.OrderClient.DeleteOrder(context.Background(), &proto.DeleteOrderRequest{
		OrderId: req.OrderID,
		UserId:  req.UserID,
	})
	if err != nil {
		http.Error(w, "Ошибка удаления заказа", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
