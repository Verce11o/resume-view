package postgres

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/pagination"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var positionLimit = 5

type PositionRepository struct {
	db *pgxpool.Pool
}

func NewPositionRepository(db *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{db: db}
}

func (p *PositionRepository) CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error) {

	id := uuid.New()
	q := "INSERT INTO positions(id, name, salary) VALUES ($1, $2, $3) RETURNING id, name, salary, created_at, updated_at"
	row, err := p.db.Query(ctx, q, id, request.Name, request.Salary)

	if err != nil {
		return models.Position{}, err
	}

	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.Position{}, err
	}

	return position, nil
}

func (p *PositionRepository) GetPosition(ctx context.Context, id string) (models.Position, error) {
	q := "SELECT id, name, salary, created_at, updated_at FROM positions WHERE id = $1"

	row, err := p.db.Query(ctx, q, id)
	if err != nil {
		return models.Position{}, err
	}

	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.Position{}, err
	}

	return position, nil
}

func (p *PositionRepository) GetPositionList(ctx context.Context, cursor string) (models.PositionList, error) {
	var createdAt time.Time
	var positionID uuid.UUID
	var err error

	if cursor != "" {
		createdAt, positionID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.PositionList{}, err
		}
	}

	q := "SELECT id, name, salary, created_at, updated_at FROM positions WHERE (created_at, id) > ($1, $2) ORDER BY created_at, id LIMIT $3"

	row, err := p.db.Query(ctx, q, createdAt, positionID, positionLimit)
	if err != nil {
		return models.PositionList{}, err
	}

	positionList, err := pgx.CollectRows(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.PositionList{}, err
	}

	var nextCursor string
	if len(positionList) > 0 {
		lastPosition := positionList[len(positionList)-1]
		nextCursor = pagination.EncodeCursor(lastPosition.CreatedAt, lastPosition.ID.String())
	}

	return models.PositionList{
		Cursor:    nextCursor,
		Positions: positionList,
	}, nil
}

func (p *PositionRepository) UpdatePosition(ctx context.Context, id string, request api.UpdatePosition) (models.Position, error) {
	q := `UPDATE positions SET name = COALESCE(NULLIF($2, ''), name), 
                     		   salary = COALESCE(NULLIF($3, 0), salary), updated_at = NOW()
                 WHERE id = $1 RETURNING id, name, salary, created_at, updated_at`

	row, err := p.db.Query(ctx, q, id, request.Name, request.Salary)
	if err != nil {
		return models.Position{}, err
	}

	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.Position{}, err
	}

	return position, nil
}

func (p *PositionRepository) DeletePosition(ctx context.Context, id string) error {

	q := "DELETE FROM positions WHERE id = $1"
	rows, err := p.db.Exec(ctx, q, id)

	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err
}
