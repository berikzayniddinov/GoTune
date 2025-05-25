package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gotune/events"
	"gotune/users/intern/entity"
	"gotune/users/intern/repository"
	"gotune/users/pkg/hash"
	"gotune/users/proto"
)

const (
	userCacheKeyPrefix   = "user:"
	userCacheExpiration  = 30 * time.Minute
	allUsersCacheKey     = "users:all"
	allUsersCacheExpTime = 5 * time.Minute
	confirmationCodeTTL  = 15 * time.Minute
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type UserService struct {
	repo           repository.UserRepository
	eventPublisher *events.EventPublisher
	cache          *redis.Client
	emailSender    EmailSender
	proto.UnimplementedUserServiceServer
}

type TestUserService struct {
	repo           repository.UserRepository
	eventPublisher *events.EventPublisher
	cache          *redis.Client
}

func NewUserService(
	repo repository.UserRepository,
	eventPublisher *events.EventPublisher,
	cache *redis.Client,
	emailSender EmailSender,
) *UserService {
	return &UserService{
		repo:           repo,
		eventPublisher: eventPublisher,
		cache:          cache,
		emailSender:    emailSender,
	}
}

func generateConfirmationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *UserService) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	session, err := s.repo.GetDatabase().Client().StartSession()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось начать сессию: %v", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := s.repo.FindByEmail(sessCtx, req.Email)
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

		if err := s.repo.Create(sessCtx, user); err != nil {
			return nil, err
		}

		code := generateConfirmationCode()
		err = s.cache.Set(ctx, "confirm_code:"+user.Email, code, confirmationCodeTTL).Err()
		if err != nil {
			log.Printf("Ошибка сохранения кода подтверждения в Redis: %v", err)
		}

		subject := "Добро пожаловать в GoTune! Подтвердите ваш email"
		body := fmt.Sprintf("Здравствуйте, %s!\n\nСпасибо за регистрацию. Ваш код подтверждения: %s\n\nКод действителен %v минут.", user.Username, code, confirmationCodeTTL.Minutes())
		if err := s.emailSender.SendEmail(user.Email, subject, body); err != nil {
			log.Printf("Ошибка отправки email: %v", err)
		}

		_ = s.eventPublisher.Publish("user_registered", map[string]string{
			"user_id": user.ID.Hex(),
			"email":   user.Email,
		})

		return user.ID.Hex(), nil
	}

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterUserResponse{
		UserId: result.(string),
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

	resp := &proto.GetUserResponse{
		UserId:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}
	data, _ := json.Marshal(resp)
	s.cache.Set(ctx, userCacheKeyPrefix+req.UserId, data, userCacheExpiration)

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

	_ = s.eventPublisher.Publish("user_deleted", map[string]string{
		"user_id": req.UserId,
		"email":   user.Email,
	})

	return &proto.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserService) DeleteAllUsersCache(ctx context.Context, req *proto.DeleteAllUsersCacheRequest) (*proto.DeleteAllUsersCacheResponse, error) {
	keys, err := s.cache.Keys(ctx, userCacheKeyPrefix+"*").Result()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get cache keys: %v", err)
	}

	keys = append(keys, allUsersCacheKey)

	if len(keys) > 0 {
		if err := s.cache.Del(ctx, keys...).Err(); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to delete cache: %v", err)
		}
	}

	return &proto.DeleteAllUsersCacheResponse{
		Success: true,
	}, nil
}

func (s *UserService) invalidateUserCache(ctx context.Context, userId string) {
	cacheKey := userCacheKeyPrefix + userId
	s.cache.Del(ctx, cacheKey)
}

func (s *UserService) ConfirmUser(ctx context.Context, req *proto.ConfirmUserRequest) (*proto.ConfirmUserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	storedCode, err := s.cache.Get(ctx, "confirm_code:"+req.Email).Result()
	if err != nil || storedCode != req.Code {
		return &proto.ConfirmUserResponse{
			Success: false,
			Message: "Invalid or expired confirmation code",
		}, nil
	}

	user.Confirmed = true
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	s.cache.Del(ctx, "confirm_code:"+req.Email)

	_ = s.eventPublisher.Publish("user_confirmed", map[string]string{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
	})

	return &proto.ConfirmUserResponse{
		Success: true,
		Message: "User confirmed successfully",
	}, nil
}
