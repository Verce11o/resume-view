package service

import (
	"context"
	"github.com/Verce11o/resume-view/internal/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ViewRepository interface {
	CreateView(ctx context.Context, resumeID, companyID string) (string, error)
	GetViews(ctx context.Context, cursor, resumeID string) (models.ViewList, error)
}

type ViewService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	repo   ViewRepository
}

func NewViewService(log *zap.SugaredLogger, tracer trace.Tracer, repo ViewRepository) *ViewService {
	return &ViewService{log: log, tracer: tracer, repo: repo}
}

func (v *ViewService) CreateView(ctx context.Context, resumeID, companyID string) (string, error) {
	ctx, span := v.tracer.Start(ctx, "viewService.CreateView")
	defer span.End()

	viewID, err := v.repo.CreateView(ctx, resumeID, companyID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	return viewID, nil
}

func (v *ViewService) GetResumeViews(ctx context.Context, cursor, resumeID string) (models.ViewList, error) {
	ctx, span := v.tracer.Start(ctx, "viewService.GetResumeViews")
	defer span.End()

	viewList, err := v.repo.GetViews(ctx, cursor, resumeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return models.ViewList{}, err
	}

	return viewList, nil
}
