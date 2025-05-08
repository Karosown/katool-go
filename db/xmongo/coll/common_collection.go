package coll

//本文件的意义在于，灵活的自定义，可以在对 mongo 进行 curd 添加打印日志、计时等功能
import (
	"context"
	"slices"

	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/ioc"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/db/xmongo/mongoutil"
	"github.com/karosown/katool-go/db/xmongo/options"
	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/xlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionFactory[T any] struct {
	*options.Client
	coll   *mongo.Collection
	logger xlog.Logger
	before func(ctx context.Context, funcName, dbName, collName string, filter *wrapper.QueryWrapper) context.Context
}

func NewCollectionFactory[T any](client *options.Client, coll *mongo.Collection, logger xlog.Logger, before func(ctx context.Context, funcName, dbName, collName string, filter *wrapper.QueryWrapper) context.Context) *CollectionFactory[T] {
	return &CollectionFactory[T]{Client: client, coll: coll, logger: logger, before: before}
}
func newCollection[T any](client *options.Client, coll *mongo.Collection, logger xlog.Logger, filter ...wrapper.QueryWrapper) *Collection[T] {
	identity := (&CollectionFactory[T]{Client: client, coll: coll, logger: logger}).Identity()
	return optional.IsTrueByFunc(cutil.IsEmpty(filter), optional.Identity(identity), func() *Collection[T] {
		if len(filter) != 1 {
			sys.Panic("the filter must be a single filter")
		}
		identity.qw = filter[0]
		return identity
	})
}

func (c *CollectionFactory[T]) Identity() *Collection[T] {
	return &Collection[T]{
		c.Client,
		c.coll,
		nil,
		c.logger,
		c.before,
	}
}

// Partition sizes[0]虚拟节点数量 sizes[1]每个虚拟节点包含的数据量大小
func (c *CollectionFactory[T]) Partition(key string, sizes ...int) *Collection[T] {
	partitionCollName := mongoutil.NewDefPartitionHelper(c.coll.Name(), sizes...).GetCollName(key)
	return ioc.GetDefFunc(partitionCollName, func() *Collection[T] {
		db := c.coll.Database()
		names, err := db.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			return newCollection[T](c.Client, db.Collection(c.coll.Name()), c.logger)
		}
		if !slices.Contains(names, partitionCollName) {
			err = db.CreateCollection(context.Background(), partitionCollName)
			// todo
			if err != nil {
				sys.Panic(err)
			}
		}
		return newCollection[T](c.Client, c.coll.Database().Collection(partitionCollName), c.logger)
	})
}
