package service

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"go.uber.org/zap"
)

type PositionRepository interface {
	CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error)
	GetPosition(ctx context.Context, id string) (models.Position, error)
	GetPositionList(ctx context.Context, cursor string) (models.PositionList, error)
	UpdatePosition(ctx context.Context, id string, request api.UpdatePosition) (models.Position, error)
	DeletePosition(ctx context.Context, id string) error
}

type PositionService struct {
	log  *zap.SugaredLogger
	repo PositionRepository
}

func NewPositionService(log *zap.SugaredLogger, repo PositionRepository) *PositionService {
	return &PositionService{log: log, repo: repo}
}

func (s *PositionService) CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error) {
	position, err := s.repo.CreatePosition(ctx, request)
	if err != nil {
		return models.Position{}, fmt.Errorf("create position: %w", err)
	}
	return position, nil
}

func (s *PositionService) GetPosition(ctx context.Context, id string) (models.Position, error) {
	position, err := s.repo.GetPosition(ctx, id)
	if err != nil {
		return models.Position{}, fmt.Errorf("get position: %w", err)
	}
	return position, nil
}

func (s *PositionService) GetPositionList(ctx context.Context, cursor string) (models.PositionList, error) {
	positionList, err := s.repo.GetPositionList(ctx, cursor)
	if err != nil {
		return models.PositionList{}, fmt.Errorf("get position list: %w", err)
	}
	return positionList, nil
}

func (s *PositionService) UpdatePosition(ctx context.Context, id string, request api.UpdatePosition) (models.Position, error) {
	position, err := s.repo.UpdatePosition(ctx, id, request)
	if err != nil {
		return models.Position{}, fmt.Errorf("update position: %w", err)
	}
	return position, nil
}

func (s *PositionService) DeletePosition(ctx context.Context, id string) error {
	err := s.repo.DeletePosition(ctx, id)
	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}
	return nil
}
