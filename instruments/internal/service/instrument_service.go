package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gotune/events"
	"gotune/instruments/internal/entity"
	"gotune/instruments/internal/repository"
	"gotune/instruments/proto"
)

type InstrumentService struct {
	repo repository.InstrumentRepository
	proto.UnimplementedInstrumentServiceServer
	eventPublisher *events.EventPublisher
}

func NewInstrumentService(repo repository.InstrumentRepository, publisher *events.EventPublisher) *InstrumentService {
	return &InstrumentService{repo: repo, eventPublisher: publisher}
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
	_ = s.eventPublisher.Publish("instrument_created", map[string]string{
		"id": id.Hex(),
	})
	return &proto.CreateInstrumentResponse{
		Id: id.Hex(),
	}, nil
}

func (s *InstrumentService) GetAllInstruments(ctx context.Context, req *proto.GetAllInstrumentsRequest) (*proto.GetAllInstrumentsResponse, error) {
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

	return &proto.GetAllInstrumentsResponse{
		Instruments: protoInstruments,
	}, nil
}

func (s *InstrumentService) DeleteInstrumentByID(ctx context.Context, req *proto.DeleteInstrumentByIDRequest) (*proto.DeleteInstrumentByIDResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return nil, err
	}

	return &proto.DeleteInstrumentByIDResponse{Success: true}, nil
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

	return &proto.UpdateInstrumentByIDResponse{Success: true}, nil
}
