package mongodb

import (
	"context"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/pagination"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const employeeLimit = 5

type EmployeeRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Database) *EmployeeRepository {
	return &EmployeeRepository{db: db, coll: db.Collection("employees")}
}

func (p *EmployeeRepository) CreateEmployee(ctx context.Context, req domain.CreateEmployee) (models.Employee, error) {
	positionColl := p.db.Collection("positions")

	callback := func(sess mongo.SessionContext) (interface{}, error) { //nolint:contextcheck
		if _, err := positionColl.InsertOne(sess, models.Position{ //nolint:contextcheck
			ID:        req.PositionID,
			Name:      req.PositionName,
			Salary:    req.Salary,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}); err != nil {
			return nil, err
		}

		if _, err := p.coll.InsertOne(sess, &models.Employee{
			ID:         req.EmployeeID,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			PositionID: req.PositionID,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := p.db.Client().StartSession()
	if err != nil {
		return models.Employee{}, err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback, txnOptions)
	if err != nil {
		return models.Employee{}, err
	}

	var employee models.Employee
	err = p.coll.FindOne(ctx, bson.M{
		"_id": req.EmployeeID,
	}).Decode(&employee)

	if err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}

func (p *EmployeeRepository) GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	var employee models.Employee

	err := p.coll.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&employee)

	if err != nil {
		return models.Employee{}, err
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
			return models.EmployeeList{}, err
		}
	}

	filter := bson.D{
		{
			Key: "$or", Value: bson.A{
				bson.M{
					"created_at": bson.M{"$gt": createdAt},
				},
				bson.M{
					"created_at": createdAt,
					"_id":        bson.M{"$gt": employeeID},
				},
			},
		},
	}

	var findOptions = options.Find()

	findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}, {Key: "_id", Value: 1}})
	findOptions.SetLimit(int64(employeeLimit))

	cur, err := p.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return models.EmployeeList{}, err
	}

	defer cur.Close(ctx)

	employees := make([]models.Employee, 0, employeeLimit)
	if err = cur.All(ctx, &employees); err != nil {
		return models.EmployeeList{}, err
	}

	var nextCursor string
	if len(employees) > 0 {
		lastEmployee := employees[len(employees)-1]
		nextCursor = pagination.EncodeCursor(lastEmployee.CreatedAt, lastEmployee.ID.String())
	}

	return models.EmployeeList{
		Cursor:    nextCursor,
		Employees: employees,
	}, nil
}

func (p *EmployeeRepository) UpdateEmployee(ctx context.Context, req domain.UpdateEmployee) (models.Employee, error) {
	positionColl := p.db.Collection("positions")

	err := positionColl.FindOne(ctx, bson.M{
		"_id": req.PositionID,
	}).Err()

	if err != nil {
		return models.Employee{}, err
	}

	filter := bson.D{{Key: "_id", Value: req.EmployeeID}}
	update := bson.D{{Key: "$set", Value: models.Employee{
		ID:         req.EmployeeID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		PositionID: req.PositionID,
		UpdatedAt:  time.Now().UTC(),
	}}}

	res := p.coll.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if res.Err() != nil {
		return models.Employee{}, res.Err()
	}

	var result models.Employee

	if err := res.Decode(&result); err != nil {
		return models.Employee{}, err
	}

	return result, nil
}

func (p *EmployeeRepository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	res, err := p.coll.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	if err != nil {
		return err
	}

	if res.DeletedCount < 1 {
		return mongo.ErrNoDocuments
	}
	return nil
}
