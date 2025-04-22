package services

import (
	"context"
	_ "errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/berik/GoTune/Domain/Entities"
	"github.com/berik/GoTune/Infrastructure/Persistence/Repositories"
	pb "github.com/berik/GoTune/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pb.UnimplementedUserServiceServer
	userRepo Repositories.UserRepository
}

func NewAuthService(userRepo Repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", req.Email)
	}

	// Проверяем, существует ли пользователь с таким username
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with username %s already exists", req.Username)
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user := &Entities.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Address:      req.Address,
		Phone:        req.Phone,
	}

	// Создаем пользователя
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.UserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *AuthService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	// Сначала получаем пользователя
	user, err := s.userRepo.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// Обновляем поля
	user.Username = req.Username
	user.Email = req.Email
	user.Address = req.Address
	user.Phone = req.Phone

	// Если указан новый пароль, хэшируем его
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Обновляем пользователя в БД
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userRepo.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *AuthService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Значение по умолчанию
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, convertUserToProto(user))
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
	}, nil
}

func convertUserToProto(user *Entities.User) *pb.User {
	return &pb.User{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Address:  user.Address,
		Phone:    user.Phone,
	}
}
