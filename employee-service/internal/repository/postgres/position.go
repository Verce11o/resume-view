package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/pagination"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const positionLimit = 20

type PositionRepository struct {
	db *pgxpool.Pool
}

func NewPositionRepository(db *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{db: db}
}

func (p *PositionRepository) CreatePosition(ctx context.Context, req domain.CreatePosition) (models.Position, error) {
	var (
		pgErr *pgconn.PgError
		rows  pgx.Rows
		err   error
	)

	q := "INSERT INTO positions(id, name, salary) VALUES ($1, $2, $3) RETURNING id, name, salary, created_at, updated_at"

	tx := extractTx(ctx)

	if tx != nil {
		rows, err = tx.Query(ctx, q, req.ID, req.Name, req.Salary)
	} else {
		rows, err = p.db.Query(ctx, q, req.ID, req.Name, req.Salary)
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("create position: %w", err)
	}

	position, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Position])

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return models.Position{}, customerrors.ErrDuplicateID
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("decode position: %w", err)
	}

	return position, nil
}

func (p *PositionRepository) GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error) {
	q := "SELECT id, name, salary, created_at, updated_at FROM positions WHERE id = $1"

	row, err := p.db.Query(ctx, q, id)
	if err != nil {
		return models.Position{}, fmt.Errorf("get position: %w", err)
	}

	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Position])

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("decode position: %w", err)
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

	row, err := p.db.Query(ctx, q, createdAt, positionID, positionLimit)
	if err != nil {
		return models.PositionList{}, fmt.Errorf("get position list: %w", err)
	}

	positionList, err := pgx.CollectRows(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.PositionList{}, fmt.Errorf("decode list: %w", err)
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

	tag, err := p.db.Exec(ctx, q, req.ID, req.Name, req.Salary)
	if err != nil {
		return models.Position{}, fmt.Errorf("update position: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	row, err := p.db.Query(ctx, `SELECT id, name, salary, created_at, updated_at FROM positions WHERE id = $1`, req.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Position{}, customerrors.ErrEmployeeNotFound
	} else if err != nil {
		return models.Position{}, fmt.Errorf("get position: %w", err)
	}

	position, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Position])

	if err != nil {
		return models.Position{}, fmt.Errorf("decode position: %w", err)
	}

	return position, nil
}

func (p *PositionRepository) DeletePosition(ctx context.Context, id uuid.UUID) error {
	q := "DELETE FROM positions WHERE id = $1"
	rows, err := p.db.Exec(ctx, q, id)

	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}

	if rows.RowsAffected() == 0 {
		return customerrors.ErrPositionNotFound
	}

	return nil
}
