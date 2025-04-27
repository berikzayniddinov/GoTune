package handler

import (
	"context"
	"encoding/json"
	"gotune/cart/proto"
	"net/http"
)

type CartHandler struct {
	CartClient proto.CartServiceClient
}

func NewCartHandler(client proto.CartServiceClient) *CartHandler {
	return &CartHandler{CartClient: client}
}

type AddToCartRequest struct {
	UserID       string `json:"user_id"`
	InstrumentID string `json:"instrument_id"`
	Quantity     int32  `json:"quantity"`
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := h.CartClient.AddToCart(context.Background(), &proto.AddToCartRequest{
		UserId:       req.UserID,
		InstrumentId: req.InstrumentID,
		Quantity:     req.Quantity,
	})
	if err != nil {
		http.Error(w, "Ошибка добавления в корзину", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Отсутствует user_id", http.StatusBadRequest)
		return
	}

	resp, err := h.CartClient.GetCart(context.Background(), &proto.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Ошибка получения корзины", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Items)
}

type RemoveFromCartRequest struct {
	UserID       string `json:"user_id"`
	InstrumentID string `json:"instrument_id"`
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var req RemoveFromCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := h.CartClient.RemoveFromCart(context.Background(), &proto.RemoveFromCartRequest{
		UserId:       req.UserID,
		InstrumentId: req.InstrumentID,
	})
	if err != nil {
		http.Error(w, "Ошибка удаления из корзины", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

type ClearCartRequest struct {
	UserID string `json:"user_id"`
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	var req ClearCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := h.CartClient.ClearCart(context.Background(), &proto.ClearCartRequest{
		UserId: req.UserID,
	})
	if err != nil {
		http.Error(w, "Ошибка очистки корзины", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
