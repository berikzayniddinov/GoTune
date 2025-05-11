package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotune/users/proto"
	"net/http"
)

type UserHandler struct {
	UserClient proto.UserServiceClient
}

func NewUserHandler(client proto.UserServiceClient) *UserHandler {
	return &UserHandler{UserClient: client}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.UserClient.RegisterUser(context.Background(), &proto.RegisterUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				http.Error(w, "Email уже зарегистрирован", http.StatusConflict)
				return
			default:
				http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"user_id": resp.UserId,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.UserClient.LoginUser(context.Background(), &proto.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "Ошибка при входе", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": resp.Token,
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := h.UserClient.GetAllUsers(context.Background(), &proto.GetAllUsersRequest{})
	if err != nil {
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp.Users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	resp, err := h.UserClient.GetUser(context.Background(), &proto.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"user_id":  resp.UserId,
		"username": resp.Username,
		"email":    resp.Email,
	})
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.UserClient.UpdateUser(context.Background(), &proto.UpdateUserRequest{
		UserId:   userID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				http.Error(w, "Пользователь не найден", http.StatusNotFound)
				return
			case codes.AlreadyExists:
				http.Error(w, "Email уже зарегистрирован другим пользователем", http.StatusConflict)
				return
			default:
				http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": resp.Success,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	resp, err := h.UserClient.DeleteUser(context.Background(), &proto.DeleteUserRequest{
		UserId: userID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				http.Error(w, "Пользователь не найден", http.StatusNotFound)
				return
			default:
				http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": resp.Success,
	})
}
func (h *UserHandler) DeleteAllUsersCache(w http.ResponseWriter, r *http.Request) {
	resp, err := h.UserClient.DeleteAllUsersCache(context.Background(), &proto.DeleteAllUsersCacheRequest{})
	if err != nil {
		http.Error(w, "Ошибка при очистке кеша", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"success": resp.Success,
	})
}
