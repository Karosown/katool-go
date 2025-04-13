package xmongo

import (
	"context"
	"fmt"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/xlog"
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

type CollectionFactoryBuilder[T any] struct {
	*mongo.Client
	DBName string
	logger xlog.Logger
}

func (m *CollectionFactoryBuilder[T]) CollName(name string) *coll.CollectionFactory[T] {
	return ioc.GetDefFunc("mongodb:"+":"+m.DBName+":"+name, func() *coll.CollectionFactory[T] {
		db := m.Client.Database(m.DBName)
		background := context.Background()
		names, err := db.ListCollectionNames(background, bson.D{})
		if err != nil {
			return coll.NewCollectionFactory[T](db.Collection(name), m.logger)
		}
		if !slices.Contains(names, name) {
			err = db.CreateCollection(background, name)
			// todo
			if err != nil {
				sys.Panic(err)
			}
		}
		return coll.NewCollectionFactory[T](db.Collection(name), m.logger)
	})
}

func (m *CollectionFactoryBuilder[T]) Transaction(ctx context.Context, fn func(stx mongo.SessionContext) (any, error)) (any, error) {
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
func (m *CollectionFactoryBuilder[T]) TransactionErr(ctx context.Context, fn func(stx mongo.SessionContext) error) error {
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

func NewCollectionFactoryBuilder[T any](DBName string, logger xlog.Logger, mc ...*mongo.Client) *CollectionFactoryBuilder[T] {
	ik := "katool:xmongdb:" + DBName
	def := ioc.GetDef(ik, mc[0])
	if cutil.IsBlank(def) {
		sys.Panic(fmt.Errorf("empty DB name: %s", ik))
	}
	return &CollectionFactoryBuilder[T]{
		def,
		ik,
		logger,
	}
}
