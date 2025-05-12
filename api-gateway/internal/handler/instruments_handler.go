package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gotune/instruments/proto"
	"net/http"
)

type InstrumentHandler struct {
	InstrumentClient proto.InstrumentServiceClient
}

func NewInstrumentHandler(client proto.InstrumentServiceClient) *InstrumentHandler {
	return &InstrumentHandler{InstrumentClient: client}
}

type CreateInstrumentRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (h *InstrumentHandler) CreateInstrument(w http.ResponseWriter, r *http.Request) {
	var req CreateInstrumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.InstrumentClient.CreateInstrument(context.Background(), &proto.CreateInstrumentRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		http.Error(w, "Ошибка при создании инструмента", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id": resp.Id,
	})
}

func (h *InstrumentHandler) GetAllInstruments(w http.ResponseWriter, r *http.Request) {
	resp, err := h.InstrumentClient.GetAllInstruments(context.Background(), &proto.GetAllInstrumentsRequest{})
	if err != nil {
		http.Error(w, "Ошибка получения инструментов", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Instruments)
}

func (h *InstrumentHandler) GetInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	resp, err := h.InstrumentClient.GetInstrumentByID(context.Background(), &proto.GetInstrumentByIDRequest{
		Id: id,
	})
	if err != nil {
		http.Error(w, "Инструмент не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *InstrumentHandler) DeleteInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := h.InstrumentClient.DeleteInstrumentByID(context.Background(), &proto.DeleteInstrumentByIDRequest{
		Id: id,
	})
	if err != nil {
		http.Error(w, "Ошибка при удалении инструмента", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

type UpdateInstrumentRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

func (h *InstrumentHandler) UpdateInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req UpdateInstrumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := h.InstrumentClient.UpdateInstrumentByID(context.Background(), &proto.UpdateInstrumentByIDRequest{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		http.Error(w, "Ошибка при обновлении инструмента", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

func (h *InstrumentHandler) ClearInstrumentCache(w http.ResponseWriter, r *http.Request) {
	// This could be a method added to the service to clear entire cache
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
