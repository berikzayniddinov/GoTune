package Controllers

import (
	"context"
	"log"

	"github.com/berik/GoTune/Internal/Services"
	"github.com/berik/GoTune/proto"
)

type AuthController struct {
	authService *Services.AuthService // Указатель на AuthService
}

func NewAuthController(authService *Services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	log.Println("Creating user:", req.Email)
	return c.authService.CreateUser(ctx, req)
}

func (c *AuthController) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	log.Println("Fetching user with ID:", req.Id)
	return c.authService.GetUser(ctx, req)
}

func (c *AuthController) GetUserByUsername(ctx context.Context, req *proto.GetUserByUsernameRequest) (*proto.UserResponse, error) {
	log.Println("Fetching user by username:", req.Username)
	return c.authService.GetUserByUsername(ctx, req)
}

func (c *AuthController) GetUserByEmail(ctx context.Context, req *proto.GetUserByEmailRequest) (*proto.UserResponse, error) {
	log.Println("Fetching user by email:", req.Email)
	return c.authService.GetUserByEmail(ctx, req)
}

func (c *AuthController) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	log.Println("Updating user:", req.Id)
	return c.authService.UpdateUser(ctx, req)
}

func (c *AuthController) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	log.Println("Deleting user:", req.Id)
	return c.authService.DeleteUser(ctx, req)
}

func (c *AuthController) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	log.Println("Listing users")
	return c.authService.ListUsers(ctx, req)
}
