package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"soya_milk_forum/app/usercenter/cmd/rpc/internal/config"
	"soya_milk_forum/app/usercenter/model"
	base_model "soya_milk_forum/common/mongo_base"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *mongo.Client
	MongoDb     *mongo.Database
	UserModel   model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, db, err := base_model.InitMongoDB(c.MongoDb.Uri, c.MongoDb.Db, c.MongoDb.MaxPoolSize)
	if err != nil {
		logx.Error(err)
		panic("Ping mongo failed")
	}

	return &ServiceContext{
		Config:      c,
		MongoClient: client,
		MongoDb:     db,
		UserModel:   model.NewUserModel(db),
	}
}
