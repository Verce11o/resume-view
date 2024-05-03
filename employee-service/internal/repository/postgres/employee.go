package postgres

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/pagination"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var employeeLimit = 5

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (p *EmployeeRepository) CreateEmployee(ctx context.Context, request api.CreateEmployee) (models.Employee, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return models.Employee{}, err
	}

	defer tx.Rollback(ctx)

	positionID := uuid.New()

	createPositionQuery := "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)"

	_, err = tx.Exec(ctx, createPositionQuery, positionID, request.PositionName, request.Salary)
	if err != nil {
		return models.Employee{}, err
	}

	employeeID := uuid.New()

	createEmployeeQuery := "INSERT INTO employees(id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, position_id,  created_at, updated_at"

	rows, err := tx.Query(ctx, createEmployeeQuery, employeeID, request.FirstName, request.LastName, positionID)

	if err != nil {
		return models.Employee{}, err
	}

	employee, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.Employee{}, err
	}

	return employee, tx.Commit(ctx)
}

func (p *EmployeeRepository) GetEmployee(ctx context.Context, id string) (models.Employee, error) {
	q := "SELECT id, first_name, last_name, position_id, created_at, updated_at FROM employees WHERE id = $1"

	row, err := p.db.Query(ctx, q, id)
	if err != nil {
		return models.Employee{}, err
	}

	employee, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}

func (p *EmployeeRepository) GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error) {
	var createdAt time.Time
	var employeeID uuid.UUID
	var err error

	if cursor != "" {
		createdAt, employeeID, err = pagination.DecodeCursor(cursor)
		if err != nil {
			return models.EmployeeList{}, err
		}
	}

	q := "SELECT id, first_name, last_name, position_id, created_at, updated_at FROM employees WHERE (created_at, id) > ($1, $2) ORDER BY created_at, id LIMIT $3"

	row, err := p.db.Query(ctx, q, createdAt, employeeID, employeeLimit)
	if err != nil {
		return models.EmployeeList{}, err
	}

	employeeList, err := pgx.CollectRows(row, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.EmployeeList{}, err
	}

	var nextCursor string
	if len(employeeList) > 0 {
		lastEmployee := employeeList[len(employeeList)-1]
		fmt.Println(lastEmployee.ID.String())
		nextCursor = pagination.EncodeCursor(lastEmployee.CreatedAt, lastEmployee.ID.String())
	}

	return models.EmployeeList{
		Cursor:    nextCursor,
		Employees: employeeList,
	}, nil
}

func (p *EmployeeRepository) UpdateEmployee(ctx context.Context, id string, request api.UpdateEmployee) (models.Employee, error) {
	q := `UPDATE employees SET first_name = COALESCE(NULLIF($2, ''), first_name), 
                     last_name = COALESCE(NULLIF($3, ''), last_name), 
                     position_id =  $4 , updated_at = NOW()
                 WHERE id = $1 RETURNING id, first_name, last_name, position_id, created_at, updated_at`

	row, err := p.db.Query(ctx, q, id, request.FirstName, request.LastName, request.PositionId)
	if err != nil {
		return models.Employee{}, err
	}

	employee, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Employee])

	if err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}

func (p *EmployeeRepository) DeleteEmployee(ctx context.Context, id string) error {

	q := "DELETE FROM employees WHERE id = $1"
	rows, err := p.db.Exec(ctx, q, id)

	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err
}
