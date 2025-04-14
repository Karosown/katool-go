package options

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Client struct {
	*mongo.Client
}

func NewClient(client *mongo.Client) *Client {
	return &Client{client}
}
func (m *Client) Transaction(ctx context.Context, fn func(stx mongo.SessionContext) (any, error)) (any, error) {
	sessionOptions := options.Session().SetDefaultReadPreference(readpref.Primary()).SetCausalConsistency(false)
	session, err := m.StartSession(sessionOptions)
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	transactionOptions := options.Transaction().SetWriteConcern(writeconcern.Majority()).SetReadConcern(readconcern.Snapshot())
	return session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		return fn(sessionCtx)
	}, transactionOptions)
}
func (m *Client) TransactionErr(ctx context.Context, fn func(stx mongo.SessionContext) error) error {
	sessionOptions := options.Session().SetDefaultReadPreference(readpref.Primary()).SetCausalConsistency(false)
	session, err := m.Client.StartSession(sessionOptions)
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	transactionOptions := options.Transaction().SetWriteConcern(writeconcern.Majority()).SetReadConcern(readconcern.Snapshot())
	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		return nil, fn(sessionCtx)
	}, transactionOptions)
	return err
}
