package services

import (
	"context"
	"fmt"

	"github.com/Verce11o/resume-view/resume-view/internal/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ViewRepository interface {
	CreateView(ctx context.Context, resumeID, companyID string) (uuid.UUID, error)
	ListResumeView(ctx context.Context, cursor, resumeID string) (models.ViewList, error)
}

type ViewService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	repo   ViewRepository
}

func NewViewService(log *zap.SugaredLogger, tracer trace.Tracer, repo ViewRepository) *ViewService {
	return &ViewService{log: log, tracer: tracer, repo: repo}
}

func (v *ViewService) CreateView(ctx context.Context, resumeID, companyID string) (uuid.UUID, error) {
	ctx, span := v.tracer.Start(ctx, "viewService.CreateView")
	defer span.End()

	viewID, err := v.repo.CreateView(ctx, resumeID, companyID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return uuid.Nil, fmt.Errorf("failed to create view: %w", err)
	}

	return viewID, nil
}

func (v *ViewService) ListResumeView(ctx context.Context, cursor, resumeID string) (models.ViewList, error) {
	ctx, span := v.tracer.Start(ctx, "viewService.ListResumeView")
	defer span.End()

	viewList, err := v.repo.ListResumeView(ctx, cursor, resumeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.ViewList{}, fmt.Errorf("failed to list resume views: %w", err)
	}

	return viewList, nil
}
