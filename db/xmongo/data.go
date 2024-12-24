package database

import (
	"context"
	"slices"

	"github.com/karosown/katool/container/ioc"
	"github.com/karosown/katool/db/xmongo/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoClient struct {
	*mongo.Client
	DBName string
}

func (m *MongoClient) CollName(name string) *coll.Collection {
	return ioc.GetDefFunc("mongodb:"+":"+m.DBName+":"+name, func() *coll.Collection {
		db := m.Client.Database(m.DBName)
		background := context.Background()
		names, err := db.ListCollectionNames(background, bson.D{})
		if err != nil {
			return coll.NewCollection(db.Collection(name))
		}
		if !slices.Contains(names, name) {
			err := db.CreateCollection(background, name)
			// todo
			if err != nil {
				panic(err)
			}
		}
		return coll.NewCollection(db.Collection(name))
	})
}

func (m *MongoClient) Transaction(ctx context.Context, fn func(stx mongo.SessionContext) (any, error)) (any, error) {
	sessionOptions := options.Session().SetDefaultReadPreference(readpref.Primary()).SetCausalConsistency(false)
	session, err := m.Client.StartSession(sessionOptions)
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	transactionOptions := options.Transaction().SetWriteConcern(writeconcern.Majority()).SetReadConcern(readconcern.Snapshot())
	return session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		return fn(sessionCtx)
	}, transactionOptions)
}
func (m *MongoClient) TransactionErr(ctx context.Context, fn func(stx mongo.SessionContext) error) error {
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
