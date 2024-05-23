package mongodb

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const positionLimit = 5

type PositionRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewPositionRepository(db *mongo.Database) *PositionRepository {
	return &PositionRepository{db: db, coll: db.Collection("positions")}
}

func (p *PositionRepository) CreatePosition(ctx context.Context, req domain.CreatePosition) (models.Position, error) {
	_, err := p.coll.InsertOne(ctx, &models.Position{
		ID:        req.ID,
		Name:      req.Name,
		Salary:    req.Salary,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.Position{}, customerrors.ErrDuplicateID
		}

		return models.Position{}, fmt.Errorf("create position: %w", err)
	}

	var position models.Position
	err = p.coll.FindOne(ctx, bson.M{
		"_id": req.ID,
	}).Decode(&position)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	if err != nil {
		return models.Position{}, fmt.Errorf("decode position: %w", err)
	}

	return position, nil
}

func (p *PositionRepository) GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error) {
	var position models.Position

	err := p.coll.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&position)

	if errors.Is(err, mongo.ErrNoDocuments) {
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
			return models.PositionList{}, fmt.Errorf("get position list: %w", err)
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
					"_id":        bson.M{"$gt": positionID},
				},
			},
		},
	}

	var findOptions = options.Find()

	findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}, {Key: "_id", Value: 1}})
	findOptions.SetLimit(int64(positionLimit))

	cur, err := p.coll.Find(ctx, filter, findOptions)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.PositionList{}, customerrors.ErrPositionNotFound
	}

	if err != nil {
		return models.PositionList{}, fmt.Errorf("find positions: %w", err)
	}

	defer cur.Close(ctx)

	positions := make([]models.Position, 0, positionLimit)
	if err = cur.All(ctx, &positions); err != nil {
		return models.PositionList{}, fmt.Errorf("find positions: %w", err)
	}

	var nextCursor string

	if len(positions) > 0 {
		lastPosition := positions[len(positions)-1]

		nextCursor = pagination.EncodeCursor(lastPosition.CreatedAt, lastPosition.ID.String())
	}

	return models.PositionList{
		Cursor:    nextCursor,
		Positions: positions,
	}, nil
}

func (p *PositionRepository) UpdatePosition(ctx context.Context, req domain.UpdatePosition) (models.Position, error) {
	filter := bson.D{{Key: "_id", Value: req.ID}}
	update := bson.D{{Key: "$set", Value: models.Position{
		ID:        req.ID,
		Name:      req.Name,
		Salary:    req.Salary,
		UpdatedAt: time.Now().UTC(),
	}}}

	res := p.coll.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return models.Position{}, customerrors.ErrPositionNotFound
	}

	if res.Err() != nil {
		return models.Position{}, fmt.Errorf("find and update position: %w", res.Err())
	}

	var result models.Position

	if err := res.Decode(&result); err != nil {
		return models.Position{}, fmt.Errorf("decode position: %w", err)
	}

	return result, nil
}

func (p *PositionRepository) DeletePosition(ctx context.Context, id uuid.UUID) error {
	res, err := p.coll.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	if err != nil {
		return fmt.Errorf("delete position: %w", err)
	}

	if res.DeletedCount < 1 {
		return customerrors.ErrPositionNotFound
	}

	return nil
}
