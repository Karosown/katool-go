package coll

//本文件的意义在于，灵活的自定义，可以在对 mongo 进行 curd 添加打印日志、计时等功能
import (
	"context"
	"slices"
	"time"

	"github.com/karosown/katool/container/ioc"
	"github.com/karosown/katool/db/pager"
	"github.com/karosown/katool/db/xmongo/mongo_util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	coll           *mongo.Collection
	collectionPool map[string]*mongo.Collection
}

func NewCollection(coll *mongo.Collection) *Collection {
	return &Collection{coll: coll}
}

// Partition sizes[0]虚拟节点数量 sizes[1]每个虚拟节点包含的数据量大小
func (c *Collection) Partition(key string, sizes ...int) *Collection {
	partitionCollName := mongo_util.NewDefPartitionHelper(c.coll.Name(), sizes...).GetCollName(key)
	return ioc.GetDefFunc(partitionCollName, func() *Collection {
		db := c.coll.Database()
		names, err := db.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			return NewCollection(db.Collection(c.coll.Name()))
		}
		if !slices.Contains(names, partitionCollName) {
			err = db.CreateCollection(context.Background(), partitionCollName)
			// todo
			if err != nil {
				panic(err)
			}
		}
		return NewCollection(c.coll.Database().Collection(partitionCollName))
	})
}
func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.coll.InsertOne(ctx, document, opts...)
}

func (c *Collection) FindOne(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	singleResult := c.coll.FindOne(ctx, filter, opts...)
	if singleResult.Err() != nil {
		return singleResult.Err()
	}
	return singleResult.Decode(result)
}

func (c *Collection) List(ctx context.Context, filter interface{}, result interface{}, opts ...*options.FindOptions) error {
	cur, err := c.coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	err = cur.All(ctx, result)
	return err
}
func (c *Collection) Count(ctx context.Context, filter interface {
}, opts ...*options.CountOptions) (int64, error) {
	return c.coll.CountDocuments(ctx, filter, opts...)

}

func (c *Collection) Page(ctx context.Context, filter interface{}, result interface{}, page *pager.Pager) error {
	documents, err := c.Count(ctx, filter)
	if err != nil {
		return err
	}
	page.Total = int(documents)
	var findoptions = &options.FindOptions{}
	if page.PageSize > 0 {
		findoptions.SetLimit(int64(page.PageSize))
		findoptions.SetSkip(int64((page.Page - 1) * page.PageSize))
		findoptions.SetSort(bson.D{{"created_at", -1}})
	}
	err = c.List(ctx, filter, result, findoptions)
	return err
}

func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.coll.UpdateOne(ctx, filter, mongo_util.StructToUpdateBSON(update, true), opts...)
}

func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.coll.DeleteOne(ctx, filter, opts...)
}
func (c *Collection) SoftDelete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	// 新增DeleteTime
	update := bson.M{
		"$set": bson.M{
			"delete_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// 使用UpdateOne而不是DeleteOne
	result, err := c.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// 将UpdateResult转换为DeleteResult
	deleteResult := &mongo.DeleteResult{
		DeletedCount: result.ModifiedCount,
	}
	return deleteResult, nil
}

// ... 添加 Find, UpdateMany, DeleteMany 等
