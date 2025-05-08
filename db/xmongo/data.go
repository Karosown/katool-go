package xmongo

import (
	"context"
	"fmt"
	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"slices"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/db/xmongo/options"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/xlog"

	"github.com/karosown/katool-go/container/ioc"
	"github.com/karosown/katool-go/db/xmongo/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionFactoryBuilder[T any] struct {
	DB     *options.Client
	DBName string
	logger xlog.Logger
	before func(ctx context.Context, funcName, dbName, collName string, filter *wrapper.QueryWrapper) (context.Context, error)
}

func (m *CollectionFactoryBuilder[T]) CollName(name string) *coll.CollectionFactory[T] {
	return ioc.GetDefFunc("mongodb:"+":"+m.DBName+":"+name, func() *coll.CollectionFactory[T] {
		db := m.DB.Database(m.DBName)
		background := context.Background()
		names, err := db.ListCollectionNames(background, bson.D{})
		if err != nil {
			return coll.NewCollectionFactory[T](m.DB, db.Collection(name), m.logger, m.before)
		}
		if !slices.Contains(names, name) {
			err = db.CreateCollection(background, name)
			// todo
			if err != nil {
				sys.Panic(err)
			}
		}
		return coll.NewCollectionFactory[T](m.DB, db.Collection(name), m.logger, m.before)
	})
}

func NewCollectionFactoryBuilder[T any](DBName string, logger xlog.Logger, before func(ctx context.Context, funcName, dbName, collName string, filter *wrapper.QueryWrapper) (context.Context, error), mc ...*mongo.Client) *CollectionFactoryBuilder[T] {
	ik := "katool:xmongdb:" + DBName
	def := ioc.GetDef(ik, mc[0])
	if cutil.IsBlank(def) {
		sys.Panic(fmt.Errorf("empty DB name: %s", ik))
	}
	return &CollectionFactoryBuilder[T]{
		&options.Client{def},
		DBName,
		logger,
		before,
	}
}
