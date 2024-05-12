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

const employeeLimit = 5

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (p *EmployeeRepository) CreateEmployee(ctx context.Context, req domain.CreateEmployee) (models.Employee, error) {
	var (
		pgErr *pgconn.PgError
		rows  pgx.Rows
		err   error
	)

	createEmployeeQuery := `INSERT INTO employees(id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4) 
                                        RETURNING id, first_name, last_name, position_id,  created_at, updated_at`

	tx := extractTx(ctx)

	if tx != nil {
		rows, err = tx.Query(ctx, createEmployeeQuery, req.EmployeeID, req.FirstName, req.LastName, req.PositionID)
	} else {
		rows, err = p.db.Query(ctx, createEmployeeQuery, req.EmployeeID, req.FirstName, req.LastName, req.PositionID)
	}

	if err != nil {
		return models.Employee{}, fmt.Errorf("inserting employee: %w", err)
	}

	employee, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Employee])

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return models.Employee{}, customerrors.ErrDuplicateID
	}

	if err != nil {
		return models.Employee{}, fmt.Errorf("decode employee: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Employee{}, fmt.Errorf("commit transaction: %w", err)
	}

	return employee, nil
}

func (p *EmployeeRepository) GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	q := "SELECT id, first_name, last_name, position_id, created_at, updated_at FROM employees WHERE id = $1"

	row, err := p.db.Query(ctx, q, id)
	if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	employee, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Employee])

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Employee{}, customerrors.ErrEmployeeNotFound
	}

	if err != nil {
		return models.Employee{}, fmt.Errorf("decode employee: %w", err)
	}

	return employee, nil
}

func (p *EmployeeRepository) GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error) {
	var (
		createdAt  time.Time
		employeeID uuid.UUID
		err        error
	)

	if cursor != "" {
		createdAt, employeeID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.EmployeeList{}, fmt.Errorf("decode cursor: %w", err)
		}
	}

	q := `SELECT id, first_name, last_name, position_id, created_at, updated_at 
		  FROM employees WHERE (created_at, id) > ($1, $2) ORDER BY created_at, id LIMIT $3`

	row, err := p.db.Query(ctx, q, createdAt, employeeID, employeeLimit)
	if err != nil {
		return models.EmployeeList{}, fmt.Errorf("get employee list: %w", err)
	}

	employeeList, err := pgx.CollectRows(row, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.EmployeeList{}, fmt.Errorf("decode employee list: %w", err)
	}

	var nextCursor string

	if len(employeeList) > 0 {
		lastEmployee := employeeList[len(employeeList)-1]
		nextCursor = pagination.EncodeCursor(lastEmployee.CreatedAt, lastEmployee.ID.String())
	}

	return models.EmployeeList{
		Cursor:    nextCursor,
		Employees: employeeList,
	}, nil
}

func (p *EmployeeRepository) UpdateEmployee(ctx context.Context, req domain.UpdateEmployee) (models.Employee, error) {
	q := `UPDATE employees
             SET first_name = COALESCE(NULLIF($2, ''), first_name),
                 last_name = COALESCE(NULLIF($3, ''), last_name),
                 position_id= COALESCE(NULLIF($4, '')::uuid, position_id)
           WHERE id = $1`

	var pgErr *pgconn.PgError

	tag, err := p.db.Exec(ctx, q, req.EmployeeID, req.FirstName, req.LastName, req.PositionID.String())

	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return models.Employee{}, customerrors.ErrPositionNotFound
	}

	if err != nil {
		return models.Employee{}, fmt.Errorf("update employee: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return models.Employee{}, customerrors.ErrEmployeeNotFound
	}

	row, err := p.db.Query(ctx, `SELECT id, first_name, last_name, position_id, created_at, updated_at 
									 FROM employees WHERE id = $1`, req.EmployeeID)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Employee{}, customerrors.ErrEmployeeNotFound
	} else if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	employee, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.Employee{}, fmt.Errorf("decode employee: %w", err)
	}

	return employee, nil
}

func (p *EmployeeRepository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	q := "DELETE FROM employees WHERE id = $1"
	rows, err := p.db.Exec(ctx, q, id)

	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}

	if rows.RowsAffected() == 0 {
		return customerrors.ErrEmployeeNotFound
	}

	return nil
}
