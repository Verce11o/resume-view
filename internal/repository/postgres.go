package repository

import (
	"context"
	"github.com/Verce11o/resume-view/internal/models"
	"github.com/Verce11o/resume-view/lib/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const paginationLimit = 20

type ViewRepository struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewViewRepository(db *pgxpool.Pool, tracer trace.Tracer) *ViewRepository {
	return &ViewRepository{db: db, tracer: tracer}
}

func (r *ViewRepository) CreateView(ctx context.Context, resumeID, companyID string) (string, error) {
	ctx, span := r.tracer.Start(ctx, "viewRepository.CreateView")
	defer span.End()

	var id string

	q := `INSERT INTO views (resume_id, company_id) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(ctx, q, resumeID, companyID).
		Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ViewRepository) GetViews(ctx context.Context, cursor, resumeID string) (models.ViewList, error) {
	ctx, span := r.tracer.Start(ctx, "viewRepository.GetViews")
	defer span.End()

	var viewedAt time.Time
	var viewID uuid.UUID
	var err error

	if cursor != "" {
		viewedAt, viewID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.ViewList{}, err
		}
	}

	var total int
	q := "SELECT COUNT(*) FROM views WHERE resume_id = $1"

	err = r.db.QueryRow(ctx, q, resumeID).Scan(&total)
	if err != nil {
		return models.ViewList{}, err
	}

	q = "SELECT id, resume_id, company_id, viewed_at FROM views WHERE (viewed_at, id) > ($1, $2) AND resume_id = $3 ORDER BY viewed_at DESC, id LIMIT $4"

	rows, err := r.db.Query(ctx, q, viewedAt, viewID, resumeID, paginationLimit)

	if err != nil {
		return models.ViewList{}, err
	}

	views, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.View])

	if err != nil {
		return models.ViewList{}, err
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
