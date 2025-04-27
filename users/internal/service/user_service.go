package service

import (
	"context"
	"fmt"
	"gotune/users/internal/entity"
	"gotune/users/internal/repository"
	"gotune/users/pkg/hash"
	"gotune/users/proto"
)

type UserService struct {
	repo repository.UserRepository
	proto.UnimplementedUserServiceServer
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
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
