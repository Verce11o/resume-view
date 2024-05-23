package pqnocache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/pagination"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const positionLimit = 5
const uniqueViolation = "23505"

type PositionRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (p *PositionRepository) CreatePosition(ctx context.Context, req domain.CreatePosition) (models.Position, error) {
	q := "INSERT INTO positions(id, name, salary) VALUES ($1, $2, $3) RETURNING id, name, salary, created_at, updated_at"

	row := p.db.QueryRowContext(ctx, q, req.ID, req.Name, req.Salary)

	var position models.Position
	err := row.Scan(&position.ID, &position.Name, &position.Salary, &position.CreatedAt, &position.UpdatedAt)

	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code == uniqueViolation {
		return models.Position{}, customerrors.ErrDuplicateID
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("create position: %w", err)
	}

	return position, nil
}

func (p *PositionRepository) GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error) {
	q := "SELECT id, name, salary, created_at, updated_at FROM positions WHERE id = $1"

	row := p.db.QueryRowContext(ctx, q, id)

	var position models.Position
	err := row.Scan(&position.ID, &position.Name, &position.Salary, &position.CreatedAt, &position.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("get position: %w", err)
	}

	return position, nil
}

func (p *PositionRepository) GetPositionList(ctx context.Context, cursor string) (models.PositionList, error) {
	var (
		createdAt  time.Time
		positionID uuid.UUID
		err        error
	)

	if cursor != "" {
		createdAt, positionID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.PositionList{}, fmt.Errorf("decode cursor: %w", err)
		}
	}

	q := `SELECT id, name, salary, created_at, updated_at FROM positions 
		  WHERE (created_at, id) > ($1, $2) ORDER BY created_at, id LIMIT $3`

	rows, err := p.db.QueryContext(ctx, q, createdAt, positionID, positionLimit)

	if err != nil {
		return models.PositionList{}, fmt.Errorf("get position list: %w", err)
	}

	if rows.Err() != nil {
		return models.PositionList{}, fmt.Errorf("get position list: %w", rows.Err())
	}

	defer rows.Close()

	var positionList []models.Position

	for rows.Next() {
		var position models.Position
		if err := rows.Scan(&position.ID, &position.Name, &position.Salary,
			&position.CreatedAt, &position.UpdatedAt); err != nil {
			return models.PositionList{}, fmt.Errorf("decode list: %w", err)
		}

		positionList = append(positionList, position)
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

func (p *PositionRepository) UpdatePosition(ctx context.Context, req domain.UpdatePosition) (models.Position, error) {
	q := `UPDATE positions SET name = COALESCE(NULLIF($2, ''), name), 
                     		   salary = COALESCE(NULLIF($3, 0), salary), updated_at = NOW()
                 WHERE id = $1`

	result, err := p.db.ExecContext(ctx, q, req.ID, req.Name, req.Salary)
	if err != nil {
		return models.Position{}, fmt.Errorf("update position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Position{}, fmt.Errorf("update position: %w", err)
	}

	if rowsAffected == 0 {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	return p.GetPosition(ctx, req.ID)
}

func (p *PositionRepository) DeletePosition(ctx context.Context, id uuid.UUID) error {
	q := "DELETE FROM positions WHERE id = $1"
	result, err := p.db.ExecContext(ctx, q, id)

	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}

	if rowsAffected == 0 {
		return customerrors.ErrPositionNotFound
	}

	return nil
}
