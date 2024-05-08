package service

import (
	"context"
	"fmt"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PositionRepository interface {
	CreatePosition(ctx context.Context, positionID uuid.UUID, request api.CreatePosition) (models.Position, error)
	GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error)
	GetPositionList(ctx context.Context, cursor string) (models.PositionList, error)
	UpdatePosition(ctx context.Context, id uuid.UUID, request api.UpdatePosition) (models.Position, error)
	DeletePosition(ctx context.Context, id uuid.UUID) error
}

type PositionCacheRepository interface {
	GetPosition(ctx context.Context, key string) (*models.Position, error)
	SetPosition(ctx context.Context, positionID string, position *models.Position) error
	DeletePosition(ctx context.Context, positionID string) error
}

type PositionService struct {
	log   *zap.SugaredLogger
	repo  PositionRepository
	cache PositionCacheRepository
}

func NewPositionService(log *zap.SugaredLogger, repo PositionRepository, cache PositionCacheRepository) *PositionService {
	return &PositionService{log: log, repo: repo, cache: cache}
}

func (s *PositionService) CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error) {
	position, err := s.repo.CreatePosition(ctx, uuid.New(), request)
	if err != nil {
		return models.Position{}, fmt.Errorf("create position: %w", err)
	}
	return position, nil
}

func (s *PositionService) GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error) {
	cachedPosition, err := s.cache.GetPosition(ctx, id.String())

	if err != nil {
		s.log.Errorf("get position from cache: %s", err)
	}

	if cachedPosition != nil {
		s.log.Debugf("returned from cache: %s", cachedPosition)
		return *cachedPosition, nil
	}

	position, err := s.repo.GetPosition(ctx, id)
	if err != nil {
		return models.Position{}, fmt.Errorf("get position: %w", err)
	}

	if err = s.cache.SetPosition(ctx, id.String(), &position); err != nil {
		s.log.Errorf("set position to cache: %s", err)
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

func (s *PositionService) UpdatePosition(ctx context.Context, id uuid.UUID, request api.UpdatePosition) (models.Position, error) {
	position, err := s.repo.UpdatePosition(ctx, id, request)
	if err != nil {
		return models.Position{}, fmt.Errorf("update position: %w", err)
	}

	if err = s.cache.DeletePosition(ctx, position.ID.String()); err != nil {
		s.log.Errorf("delete position from cache: %s", err)
	}

	return position, nil
}

func (s *PositionService) DeletePosition(ctx context.Context, id uuid.UUID) error {
	position, err := s.repo.GetPosition(ctx, id)
	if err != nil {
		return fmt.Errorf("get position: %w", err)
	}

	err = s.repo.DeletePosition(ctx, id)
	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}

	if err := s.cache.DeletePosition(ctx, position.ID.String()); err != nil {
		s.log.Errorf("delete position from cache: %s", err)
	}
	return nil
}
