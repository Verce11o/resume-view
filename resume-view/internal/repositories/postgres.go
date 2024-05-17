package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	customerrors "github.com/Verce11o/resume-view/resume-view/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/resume-view/internal/lib/pagination"
	"github.com/Verce11o/resume-view/resume-view/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

const paginationLimit = 20

type ViewRepository struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewViewRepository(db *pgxpool.Pool, tracer trace.Tracer) *ViewRepository {
	return &ViewRepository{db: db, tracer: tracer}
}

func (r *ViewRepository) CreateView(ctx context.Context, resumeID, companyID string) (uuid.UUID, error) {
	ctx, span := r.tracer.Start(ctx, "viewRepository.CreateView")
	defer span.End()

	var id uuid.UUID

	q := `INSERT INTO views (resume_id, company_id) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(ctx, q, resumeID, companyID).
		Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed create view: %w", err)
	}

	return id, nil
}

func (r *ViewRepository) ListResumeView(ctx context.Context, cursor, resumeID string) (models.ViewList, error) {
	ctx, span := r.tracer.Start(ctx, "viewRepository.ListResumeView")
	defer span.End()

	var (
		viewedAt time.Time
		viewID   uuid.UUID
		err      error
	)

	if cursor != "" {
		viewedAt, viewID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.ViewList{}, customerrors.ErrInvalidCursor
		}
	}

	var total int

	q := "SELECT COUNT(*) FROM views WHERE resume_id = $1"

	err = r.db.QueryRow(ctx, q, resumeID).Scan(&total)
	if err != nil && errors.Is(err, pgx.ErrNoRows) || total == 0 {
		return models.ViewList{}, customerrors.ErrNotFound
	}

	if err != nil {
		return models.ViewList{}, fmt.Errorf("failed to count views: %w", err)
	}

	q = `SELECT id, resume_id, company_id, viewed_at FROM views WHERE (viewed_at, id) > ($1, $2) 
		 AND resume_id = $3 ORDER BY viewed_at DESC, id LIMIT $4`

	rows, err := r.db.Query(ctx, q, viewedAt, viewID, resumeID, paginationLimit)

	if err != nil && errors.Is(err, pgx.ErrNoRows) || total == 0 {
		return models.ViewList{}, customerrors.ErrNotFound
	}

	if err != nil {
		return models.ViewList{}, fmt.Errorf("failed to list views: %w", err)
	}

	views, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.View])

	if err != nil {
		return models.ViewList{}, fmt.Errorf("failed to list views: %w", err)
	}

	var nextCursor string
	if len(views) > 0 {
		nextCursor = pagination.EncodeCursor(views[len(views)-1].ViewedAt, views[len(views)-1].ID.String())
	}

	return models.ViewList{
		Cursor: nextCursor,
		Views:  views,
		Total:  total,
	}, nil
}
