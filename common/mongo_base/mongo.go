package base_model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// InitMongoDB 初始化MongoDB
func InitMongoDB(uri, dbName string, maxPoolSize uint64) (*mongo.Client, *mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(maxPoolSize)
	// 连接MongoDB
	client, err := mongo.Connect(ctx, dbOptions)
	if err != nil {
		return nil, nil, err
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(dbName)
	if err != nil {
		return nil, nil, err
	}

	return client, db, nil
}

var dbPool *mongo.Database

func GetDB() *mongo.Database {
	return dbPool
}

func GetCollection(coName string) *mongo.Collection {
	return GetDB().Collection(coName)
}

func SetDB(dataBase *mongo.Database) {
	dbPool = dataBase
}
