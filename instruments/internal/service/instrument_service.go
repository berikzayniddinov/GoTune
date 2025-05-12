package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gotune/events"
	"gotune/instruments/internal/entity"
	"gotune/instruments/internal/repository"
	"gotune/instruments/proto"
)

const (
	instrumentCacheKeyPrefix   = "instrument:"
	allInstrumentsCacheKey     = "instruments:all"
	instrumentCacheExpTime     = 30 * time.Minute
	allInstrumentsCacheExpTime = 5 * time.Minute
)

type InstrumentService struct {
	repo           repository.InstrumentRepository
	eventPublisher *events.EventPublisher
	cache          *redis.Client
	proto.UnimplementedInstrumentServiceServer
}

func NewInstrumentService(repo repository.InstrumentRepository, publisher *events.EventPublisher, cache *redis.Client) *InstrumentService {
	return &InstrumentService{
		repo:           repo,
		eventPublisher: publisher,
		cache:          cache,
	}
}

func (s *InstrumentService) CreateInstrument(ctx context.Context, req *proto.CreateInstrumentRequest) (*proto.CreateInstrumentResponse, error) {
	instrument := &entity.Instrument{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
	id, err := s.repo.Create(ctx, instrument)
	if err != nil {
		return nil, err
	}

	// Invalidate global instruments list cache
	s.cache.Del(ctx, allInstrumentsCacheKey)

	_ = s.eventPublisher.Publish("instrument_created", map[string]string{
		"id": id.Hex(),
	})
	return &proto.CreateInstrumentResponse{
		Id: id.Hex(),
	}, nil
}

func (s *InstrumentService) GetInstrumentByID(ctx context.Context, req *proto.GetInstrumentByIDRequest) (*proto.Instrument, error) {
	cacheKey := instrumentCacheKeyPrefix + req.Id
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResp proto.Instrument
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	inst, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &proto.Instrument{
		Id:          inst.ID.Hex(),
		Name:        inst.Name,
		Description: inst.Description,
		Price:       inst.Price,
	}

	data, _ := json.Marshal(resp)
	s.cache.Set(ctx, cacheKey, data, instrumentCacheExpTime)

	return resp, nil
}

func (s *InstrumentService) GetAllInstruments(ctx context.Context, req *proto.GetAllInstrumentsRequest) (*proto.GetAllInstrumentsResponse, error) {
	cached, err := s.cache.Get(ctx, allInstrumentsCacheKey).Result()
	if err == nil {
		var cachedResp proto.GetAllInstrumentsResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	instruments, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var protoInstruments []*proto.Instrument
	for _, inst := range instruments {
		protoInstruments = append(protoInstruments, &proto.Instrument{
			Id:          inst.ID.Hex(),
			Name:        inst.Name,
			Description: inst.Description,
			Price:       inst.Price,
		})
	}

	resp := &proto.GetAllInstrumentsResponse{
		Instruments: protoInstruments,
	}

	data, _ := json.Marshal(resp)
	s.cache.Set(ctx, allInstrumentsCacheKey, data, allInstrumentsCacheExpTime)

	return resp, nil
}

func (s *InstrumentService) UpdateInstrumentByID(ctx context.Context, req *proto.UpdateInstrumentByIDRequest) (*proto.UpdateInstrumentByIDResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	instrument := &entity.Instrument{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.repo.UpdateByID(ctx, id, instrument); err != nil {
		return nil, err
	}

	s.cache.Del(ctx, instrumentCacheKeyPrefix+req.Id)
	s.cache.Del(ctx, allInstrumentsCacheKey)

	_ = s.eventPublisher.Publish("instrument_updated", map[string]string{
		"id": req.Id,
	})

	return &proto.UpdateInstrumentByIDResponse{Success: true}, nil
}

func (s *InstrumentService) DeleteInstrumentByID(ctx context.Context, req *proto.DeleteInstrumentByIDRequest) (*proto.DeleteInstrumentByIDResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return nil, err
	}

	s.cache.Del(ctx, instrumentCacheKeyPrefix+req.Id)
	s.cache.Del(ctx, allInstrumentsCacheKey)

	_ = s.eventPublisher.Publish("instrument_deleted", map[string]string{
		"id": req.Id,
	})

	return &proto.DeleteInstrumentByIDResponse{Success: true}, nil
}
