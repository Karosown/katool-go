package coll

import (
	"context"
	"time"

	"github.com/karosown/katool-go/db/pager"
	"github.com/karosown/katool-go/db/xmongo/mongoutil"
	options2 "github.com/karosown/katool-go/db/xmongo/options"
	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"github.com/karosown/katool-go/xlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[T any] struct {
	*options2.Client
	coll   *mongo.Collection
	qw     wrapper.QueryWrapper
	logger xlog.Logger
	before func(ctx context.Context, funcName, dbName, collName string, filter *wrapper.QueryWrapper, entity *T) (context.Context, error)
}

func (c *Collection[T]) Coll() *mongo.Collection {
	return c.coll
}

func (c *Collection[T]) Query(filter wrapper.QueryWrapper) *Collection[T] {
	if c.logger != nil {
		c.logger.Info("MongoDB/DocumentDB Query Bson is {}", filter.ToJSON())
	}
	return newCollection[T](c.Client, c.coll, c.logger, filter)
}

func (c *Collection[T]) InsertOne(ctx context.Context, document *T, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if c.before != nil {
		transaction, err := c.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
			ctx, err := c.before(stx, "InsertOne", c.coll.Database().Name(), c.coll.Name(), &c.qw, document)
			if err != nil {
				return ctx, err
			}
			return c.coll.InsertOne(ctx, document, opts...)
		})
		return transaction.(*mongo.InsertOneResult), err
	} else {
		return c.coll.InsertOne(ctx, document, opts...)
	}
}

func (c *Collection[T]) FindOne(ctx context.Context, result *T, opts ...*options.FindOneOptions) error {
	singleResult := c.coll.FindOne(ctx, c.filter(), opts...)
	if singleResult.Err() != nil {
		return singleResult.Err()
	}
	return singleResult.Decode(result)
}

func (c *Collection[T]) List(ctx context.Context, opts ...*options.FindOptions) (*[]T, error) {
	result := &[]T{}
	if c.before != nil {
		return nil, c.TransactionErr(ctx, func(stx mongo.SessionContext) error {
			ctx, err := c.before(stx, "List", c.coll.Database().Name(), c.coll.Name(), &c.qw, nil)
			if err != nil {
				return err
			}
			cur, err := c.coll.Find(ctx, c.filter(), opts...)
			if err != nil {
				return err
			}
			return cur.All(ctx, result)
		})
	}
	cur, err := c.coll.Find(ctx, c.filter(), opts...)
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, result)
	return result, err
}
func (c *Collection[T]) Count(ctx context.Context, opts ...*options.CountOptions) (int64, error) {
	if c.before != nil {
		transaction, err := c.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
			ctx, err := c.before(stx, "Count", c.coll.Database().Name(), c.coll.Name(), &c.qw, nil)
			if err != nil {
				return ctx, err
			}
			return c.coll.CountDocuments(ctx, c.filter(), opts...)
		})
		return transaction.(int64), err
	}
	return c.coll.CountDocuments(ctx, c.filter(), opts...)
}

func (c *Collection[T]) Page(ctx context.Context, page *pager.Pager) (*[]T, error) {
	documents, err := c.Count(ctx)
	if err != nil {
		return nil, err
	}
	page.Total = int(documents)
	var findoptions = &options.FindOptions{}
	if page.PageSize > 0 {
		findoptions.SetLimit(int64(page.PageSize))
		findoptions.SetSkip(int64((page.Page - 1) * page.PageSize))
		findoptions.SetSort(bson.D{{"created_at", -1}})
	}
	return c.List(ctx, findoptions)
}

func (c *Collection[T]) UpdateOne(ctx context.Context, update *T, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if c.before != nil {
		transaction, err := c.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
			ctx, err := c.before(stx, "UpdateOne", c.coll.Database().Name(), c.coll.Name(), &c.qw, update)
			if err != nil {
				return ctx, err
			}
			return c.coll.UpdateOne(ctx, c.filter(), mongoutil.StructToUpdateBSON(update, true), opts...)
		})
		return transaction.(*mongo.UpdateResult), err
	} else {
		return c.coll.UpdateOne(ctx, c.filter(), mongoutil.StructToUpdateBSON(update, true), opts...)
	}
}

func (c *Collection[T]) DeleteOne(ctx context.Context, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if c.before != nil {
		transaction, err := c.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
			ctx, err := c.before(ctx, "DeleteOne", c.coll.Database().Name(), c.coll.Name(), &c.qw, nil)
			if err != nil {
				return ctx, err
			}
			return c.coll.DeleteOne(ctx, c.filter(), opts...)
		})
		return transaction.(*mongo.DeleteResult), err
	}
	return c.coll.DeleteOne(ctx, c.filter(), opts...)
}
func (c *Collection[T]) SoftDelete(ctx context.Context, opts ...*options.UpdateOptions) (*mongo.DeleteResult, error) {
	// 新增DeleteTime
	update := bson.M{
		"$set": bson.M{
			wrapper.DeletedField: primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	if c.before != nil {
		transaction, err := c.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
			ctx, err := c.before(ctx, "SoftDelete", c.coll.Database().Name(), c.coll.Name(), &c.qw, nil)
			if err != nil {
				return ctx, err
			}
			// 使用UpdateOne而不是DeleteOne
			result, err := c.coll.UpdateMany(ctx, c.filter(), update, opts...)
			if err != nil {
				return nil, err
			}

			// 将UpdateResult转换为DeleteResult
			deleteResult := &mongo.DeleteResult{
				DeletedCount: result.ModifiedCount,
			}
			return deleteResult, nil
		})
		return transaction.(*mongo.DeleteResult), err
	}
	// 使用UpdateOne而不是DeleteOne
	result, err := c.coll.UpdateMany(ctx, c.filter(), update, opts...)
	if err != nil {
		return nil, err
	}

	// 将UpdateResult转换为DeleteResult
	deleteResult := &mongo.DeleteResult{
		DeletedCount: result.ModifiedCount,
	}
	return deleteResult, nil
}

func (c *Collection[T]) filter() wrapper.QueryWrapper {
	return c.qw
}

// ... 添加 Find, UpdateMany, DeleteMany 等
