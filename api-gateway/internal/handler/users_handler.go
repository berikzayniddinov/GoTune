package handler

import (
	"context"
	"encoding/json"
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
