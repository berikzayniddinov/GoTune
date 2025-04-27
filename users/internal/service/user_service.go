package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotune/events"
	"gotune/users/internal/entity"
	"gotune/users/internal/repository"
	"gotune/users/pkg/hash"
	"gotune/users/proto"
)

type UserService struct {
	repo           repository.UserRepository
	eventPublisher *events.EventPublisher
	proto.UnimplementedUserServiceServer
}

func NewUserService(repo repository.UserRepository, eventPublisher *events.EventPublisher) *UserService {
	return &UserService{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Email уже зарегистрирован")
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	_ = s.eventPublisher.Publish("user_registered", map[string]string{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
	})

	return &proto.RegisterUserResponse{
		UserId: user.ID.Hex(),
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("неправильный пароль")
	}

	token := "fake-jwt-token-for-" + user.ID.Hex()

	return &proto.LoginUserResponse{
		Token: token,
	}, nil
}
func (s *UserService) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.User
	for _, u := range users {
		protoUsers = append(protoUsers, &proto.User{
			Id:        u.ID.Hex(),
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	return &proto.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.GetUserResponse{
		UserId:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
