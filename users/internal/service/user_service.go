package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotune/events"
	"gotune/users/internal/entity"
	"gotune/users/internal/repository"
	"gotune/users/pkg/hash"
	"gotune/users/proto"
	"time"
)

const (
	userCacheKeyPrefix   = "user:"
	userCacheExpiration  = 30 * time.Minute
	allUsersCacheKey     = "users:all"
	allUsersCacheExpTime = 5 * time.Minute
)

type UserService struct {
	repo           repository.UserRepository
	eventPublisher *events.EventPublisher
	cache          *redis.Client
	proto.UnimplementedUserServiceServer
}

func NewUserService(repo repository.UserRepository, eventPublisher *events.EventPublisher, cache *redis.Client) *UserService {
	return &UserService{
		repo:           repo,
		eventPublisher: eventPublisher,
		cache:          cache,
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

	//s.invalidateUserCache(ctx, user.ID.Hex())

	//s.cache.Del(ctx, allUsersCacheKey)

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
	cached, err := s.cache.Get(ctx, allUsersCacheKey).Result()
	if err == nil {
		var cachedResp proto.GetAllUsersResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

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

	resp := &proto.GetAllUsersResponse{
		Users: protoUsers,
	}

	data, _ := json.Marshal(resp)
	s.cache.Set(ctx, allUsersCacheKey, data, allUsersCacheExpTime)

	return resp, nil
}

func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	cacheKey := userCacheKeyPrefix + req.UserId
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResp proto.GetUserResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	resp := &proto.GetUserResponse{
		UserId:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}

	data, _ := json.Marshal(resp)
	s.cache.Set(ctx, cacheKey, data, userCacheExpiration)

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		existingUser, err := s.repo.FindByEmail(ctx, req.Email)
		if err == nil && existingUser.ID.Hex() != req.UserId {
			return nil, status.Errorf(codes.AlreadyExists, "Email already registered")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := hash.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx, req.UserId)

	s.cache.Del(ctx, allUsersCacheKey)

	_ = s.eventPublisher.Publish("user_updated", map[string]string{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
	})

	return &proto.UpdateUserResponse{
		Success: true,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if err := s.repo.Delete(ctx, req.UserId); err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx, req.UserId)

	s.cache.Del(ctx, allUsersCacheKey)

	_ = s.eventPublisher.Publish("user_deleted", map[string]string{
		"user_id": req.UserId,
		"email":   user.Email,
	})

	return &proto.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserService) invalidateUserCache(ctx context.Context, userId string) {
	cacheKey := userCacheKeyPrefix + userId
	s.cache.Del(ctx, cacheKey)
}
