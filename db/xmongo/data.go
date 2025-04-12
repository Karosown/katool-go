package xmongo

import (
	"context"
	"fmt"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/sys"
	"slices"

	"github.com/karosown/katool-go/container/ioc"
	"github.com/karosown/katool-go/db/xmongo/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoClient[T any] struct {
	*mongo.Client
	DBName string
}

func (m *MongoClient[T]) CollName(name string) *coll.Collection[T] {
	return ioc.GetDefFunc("mongodb:"+":"+m.DBName+":"+name, func() *coll.Collection[T] {
		db := m.Client.Database(m.DBName)
		background := context.Background()
		names, err := db.ListCollectionNames(background, bson.D{})
		if err != nil {
			return coll.NewCollection[T](db.Collection(name))
		}
		if !slices.Contains(names, name) {
			err = db.CreateCollection(background, name)
			// todo
			if err != nil {
				sys.Panic(err)
			}
		}
		return coll.NewCollection[T](db.Collection(name))
	})
}

func (m *MongoClient[T]) Transaction(ctx context.Context, fn func(stx mongo.SessionContext) (any, error)) (any, error) {
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
func (m *MongoClient[T]) TransactionErr(ctx context.Context, fn func(stx mongo.SessionContext) error) error {
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

func NewMongoClient[T any](DBName string, mc ...*mongo.Client) *MongoClient[T] {
	ik := "katool:xmongdb:" + DBName
	def := ioc.GetDef(ik, mc[0])
	if cutil.IsBlank(def) {
		sys.Panic(fmt.Errorf("empty DB name: %s", ik))
	}
	return &MongoClient[T]{
		def,
		ik,
	}
}
