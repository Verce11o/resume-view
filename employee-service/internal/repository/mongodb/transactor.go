package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Transactor struct {
	client *mongo.Client
}

func NewTransactor(client *mongo.Client) *Transactor {
	return &Transactor{client: client}
}

func (t *Transactor) WithTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := t.client.StartSession()
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) { //nolint:contextcheck
		return nil, tFunc(sessCtx)
	}, txnOptions)

	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	return nil
}
