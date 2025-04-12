package coll

//本文件的意义在于，灵活的自定义，可以在对 mongo 进行 curd 添加打印日志、计时等功能
import (
	"context"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/ioc"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/db/xmongo/mongo_util"
	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"github.com/karosown/katool-go/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"slices"
)

type CommonCollection[T any] struct {
	coll *mongo.Collection
}

func NewCollection[T any](coll *mongo.Collection, filter ...wrapper.QueryWrapper) *Collection[T] {
	identity := (&CommonCollection[T]{coll: coll}).Identity()
	return optional.IsTrueByFunc(cutil.IsEmpty(filter), optional.Identity(identity), func() *Collection[T] {
		if len(filter) != 1 {
			sys.Panic("the filter must be a single filter")
		}
		identity.qw = filter[0]
		return identity
	})
}

func (c *CommonCollection[T]) Identity() *Collection[T] {
	return &Collection[T]{
		c.coll,
		nil,
	}
}

// Partition sizes[0]虚拟节点数量 sizes[1]每个虚拟节点包含的数据量大小
func (c *CommonCollection[T]) Partition(key string, sizes ...int) *Collection[T] {
	partitionCollName := mongo_util.NewDefPartitionHelper(c.coll.Name(), sizes...).GetCollName(key)
	return ioc.GetDefFunc(partitionCollName, func() *Collection[T] {
		db := c.coll.Database()
		names, err := db.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			return NewCollection[T](db.Collection(c.coll.Name()))
		}
		if !slices.Contains(names, partitionCollName) {
			err = db.CreateCollection(context.Background(), partitionCollName)
			// todo
			if err != nil {
				sys.Panic(err)
			}
		}
		return NewCollection[T](c.coll.Database().Collection(partitionCollName))
	})
}
