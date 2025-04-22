package Interfaces

import (
	"context"
	pb "github.com/berik/GoTune/proto"
)

type AuthService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error)
	GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error)
	GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
}
