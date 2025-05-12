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

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type AddToCartRequest struct {
	UserID       string `json:"user_id"`
	InstrumentID string `json:"instrument_id"`
	Quantity     int32  `json:"quantity"`
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "Пустое тело запроса")
		return
	}
	defer r.Body.Close()

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	_, err := h.CartClient.AddToCart(context.Background(), &proto.AddToCartRequest{
		UserId:       req.UserID,
		InstrumentId: req.InstrumentID,
		Quantity:     req.Quantity,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка добавления в корзину")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Отсутствует user_id")
		return
	}

	resp, err := h.CartClient.GetCart(context.Background(), &proto.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка получения корзины")
		return
	}

	writeJSON(w, http.StatusOK, resp.Items)
}

type RemoveFromCartRequest struct {
	UserID       string `json:"user_id"`
	InstrumentID string `json:"instrument_id"`
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "Пустое тело запроса")
		return
	}
	defer r.Body.Close()

	var req RemoveFromCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	_, err := h.CartClient.RemoveFromCart(context.Background(), &proto.RemoveFromCartRequest{
		UserId:       req.UserID,
		InstrumentId: req.InstrumentID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка удаления из корзины")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

type ClearCartRequest struct {
	UserID string `json:"user_id"`
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeError(w, http.StatusBadRequest, "Пустое тело запроса")
		return
	}
	defer r.Body.Close()

	var req ClearCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	_, err := h.CartClient.ClearCart(context.Background(), &proto.ClearCartRequest{
		UserId: req.UserID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка очистки корзины")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
