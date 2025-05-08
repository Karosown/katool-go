package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/karosown/katool-go/db/xmongo/wrapper"
	"github.com/karosown/katool-go/sys"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/karosown/katool-go/db/xmongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test data structure
type UserInfo struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

func init() {
	// 定义软删除字段，默认是delete_at
	wrapper.DeletedField = "delete_at"
	// 这个构造器会加上mongo_util.BaseFilter
	wrapper.NewQuery().Eq("_id", "1").Build()
	// 如果想要构建原始查询，使用
	wrapper.NewQuery().Eq("_id", "1").Origin()
}
func setupTestMongoClient(t *testing.T) *xmongo.CollectionFactoryBuilder[UserInfo] {
	// Connect to a test MongoDB (preferably use a test database or Docker container)
	// For testing purposes, we connect to a local MongoDB instance
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	// Create a new client with a test database name
	return xmongo.NewCollectionFactoryBuilder[UserInfo]("test_db", nil, nil, client)
}

func cleanupTestMongoClient(t *testing.T, client *xmongo.CollectionFactoryBuilder[UserInfo]) {
	ctx := context.Background()
	err := client.DB.Client.Database("test_db").Drop(ctx)
	assert.NoError(t, err)

	err = client.DB.Client.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestMongoClientCollName(t *testing.T) {
	client := setupTestMongoClient(t)
	defer cleanupTestMongoClient(t, client)

	// Test getting a collection
	collection := client.CollName("test_collection")
	assert.NotNil(t, collection)
}

func TestMongoClientTransaction(t *testing.T) {
	client := setupTestMongoClient(t)
	defer cleanupTestMongoClient(t, client)

	ctx := context.Background()

	// Test a successful transaction
	result, err := client.DB.Transaction(ctx, func(stx mongo.SessionContext) (any, error) {
		// Do some operations inside the transaction
		collection := client.CollName("test_collection").Identity()
		_, err := collection.InsertOne(stx, &UserInfo{ID: primitive.NewObjectID(), Name: "Test", Age: 30})
		return "success", err
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)

	// Verify the data was inserted
	collection := client.CollName("test_collection").Identity()
	var data UserInfo
	err = collection.Query(wrapper.NewQuery().Eq("_id", "1").Build()).FindOne(ctx, &data)
	assert.NoError(t, err)
	assert.Equal(t, "Test", data.Name)
}

func TestMongoClientTransactionErr(t *testing.T) {
	client := setupTestMongoClient(t)
	defer cleanupTestMongoClient(t, client)

	ctx := context.Background()

	// Test transaction with error handling
	err := client.DB.TransactionErr(ctx, func(stx mongo.SessionContext) error {
		collection := client.CollName("test_collection").Identity()
		_, err := collection.InsertOne(stx, &UserInfo{ID: primitive.NewObjectID(), Name: "Another Test", Age: 25})
		return err
	})

	assert.NoError(t, err)

	// Verify the data was inserted
	collection := client.CollName("test_collection").Identity()
	var data UserInfo
	err = collection.Query(wrapper.NewQuery().Eq("_id", "2").Build()).FindOne(ctx, &data)
	assert.NoError(t, err)
	assert.Equal(t, "Another Test", data.Name)
}

func TestNewMongoClient(t *testing.T) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	// Test creating a new client
	mongoClient := xmongo.NewCollectionFactoryBuilder[UserInfo]("test_db", nil, nil, client)
	assert.NotNil(t, mongoClient)
	assert.Equal(t, "katool:xmongdb:test_db", mongoClient.DBName)
}
func TestAll(t *testing.T) {
	// 连接到MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetDirect(true))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	// 创建MongoClient
	mongoClient := xmongo.NewCollectionFactoryBuilder[UserInfo]("eshop", nil, nil, client)

	// 获取集合
	collection := mongoClient.CollName("userinfos").Identity()

	// 插入数据
	userData := UserInfo{
		ID:   primitive.NewObjectID(),
		Name: "张三1",
		Age:  30,
	}

	//使用事务插入数据
	//_, err = mongoClient.Transaction(ctx, func(stx mongo.SessionContext) (interface{}, error) {
	_, err = collection.InsertOne(ctx, &userData)
	if err != nil {
		sys.Panic("插入数据失败: " + err.Error())
	}
	//return nil, nil
	//})

	//if err != nil {
	//	log.Fatalf("事务执行失败: %v", err)
	//}

	// 查询数据
	var result UserInfo
	err = collection.
		Query(wrapper.NewQuery().
			Eq("name", "张三1").
			Build()).
		FindOne(ctx, &result)
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}

	fmt.Printf("查询结果: ID=%s, Name=%s, Age=%d\n", result.ID, result.Name, result.Age)

	// 更新数据
	//err = mongoClient.TransactionErr(ctx, func(stx mongo.SessionContext) error {
	_, err = collection.Query(map[string]interface{}{"_id": result.ID}).UpdateOne(
		ctx,
		&UserInfo{Age: 78},
		nil,
	)
	//return err
	//})

	if err != nil {
		log.Fatalf("更新数据失败: %v", err)
	}

	// 再次查询确认更新
	err = collection.Query(wrapper.QueryWrapper{"_id": result.ID}).FindOne(ctx, &result)
	if err != nil {
		log.Fatalf("查询更新后数据失败: %v", err)
	}

	fmt.Printf("更新后结果: ID=%s, Name=%s, Age=%d\n", result.ID, result.Name, result.Age)

	// 删除数据
	//err = mongoClient.TransactionErr(ctx, func(stx mongo.SessionContext) error {
	_, err = collection.Query(wrapper.QueryWrapper{"_id": result.ID}).DeleteOne(ctx)
	//return err
	//})

	if err != nil {
		log.Fatalf("删除数据失败: %v", err)
	}

	fmt.Println("示例完成")
}

func TestTransaction(t *testing.T) {
	// 连接到MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetDirect(true))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	// 创建MongoClient
	mongoClient := xmongo.NewCollectionFactoryBuilder[UserInfo]("eshop", nil, nil, client)

	// 获取集合
	collection := mongoClient.CollName("userinfos").Identity()

	// 插入数据
	userData := UserInfo{
		ID:   primitive.NewObjectID(),
		Name: "张三1",
		Age:  30,
	}

	//使用事务插入数据
	_, err = collection.Transaction(ctx, func(stx mongo.SessionContext) (interface{}, error) {
		_, err = collection.InsertOne(stx, &userData)
		if err != nil {
			sys.Panic("插入数据失败: " + err.Error())
		}
		return nil, nil
	})

	if err != nil {
		log.Fatalf("事务执行失败: %v", err)
	}

	// 查询数据
	var result UserInfo
	err = collection.
		Query(wrapper.NewQuery().
			Eq("name", "张三1").
			Build()).
		FindOne(ctx, &result)
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}

	fmt.Printf("查询结果: ID=%s, Name=%s, Age=%d\n", result.ID, result.Name, result.Age)

	// 更新数据
	err = collection.TransactionErr(ctx, func(stx mongo.SessionContext) error {
		_, err = collection.Query(map[string]interface{}{"_id": result.ID}).UpdateOne(
			stx,
			&UserInfo{Age: 78},
			nil,
		)
		return err
	})

	if err != nil {
		log.Fatalf("更新数据失败: %v", err)
	}

	// 再次查询确认更新
	err = collection.Query(wrapper.QueryWrapper{"_id": result.ID}).FindOne(ctx, &result)
	if err != nil {
		log.Fatalf("查询更新后数据失败: %v", err)
	}

	fmt.Printf("更新后结果: ID=%s, Name=%s, Age=%d\n", result.ID, result.Name, result.Age)

	// 删除数据
	err = collection.TransactionErr(ctx, func(stx mongo.SessionContext) error {
		_, err = collection.Query(wrapper.QueryWrapper{"_id": result.ID}).DeleteOne(stx)
		return err
	})

	if err != nil {
		log.Fatalf("删除数据失败: %v", err)
	}

	fmt.Println("示例完成")
}
