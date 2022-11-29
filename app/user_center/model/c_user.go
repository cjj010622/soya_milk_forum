package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"soya_milk_forum/common/mongo_base"
)

type (
	User struct {
		Account             string `bson:"account"`
		Password            string `bson:"password"`
		Email               string `bson:"email"`
		TelephoneNumber     string `bson:"telephone_number"`
		Status              int8   `bson:"status"`
		Avatar              string `bson:"avatar"`
		Type                int8
		base_model.BaseBean `bson:",inline"`
	}

	UserModel interface {
		base_model.Interface
	}

	UserModelImpl struct {
		db        *mongo.Database
		tableName string
	}
)

func (u *UserModelImpl) TableName() string {
	return u.tableName
}

func (u *UserModelImpl) Collection() *mongo.Collection {
	return u.db.Collection(u.tableName)
}

func (u *UserModelImpl) DelMany(ctx context.Context, ids []primitive.ObjectID, deletedBy int64) error {
	return base_model.DefaultDelMany(ctx, u.Collection(), "_id", ids, deletedBy)
}

func (u *UserModelImpl) Add(ctx context.Context, user User) (id string, err error) {
	co := u.Collection()
	//初始化时间
	base_model.InitTime(user)
	var insertId *mongo.InsertOneResult
	insertId, err = co.InsertOne(ctx, user, nil)

	id = insertId.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (u *UserModelImpl) Update(ctx context.Context, user User) (err error) {
	update := base_model.StructureBsonD(user, base_model.MODEL2)

	_, err = base_model.UpsertOne(ctx, u.Collection(), base_model.DefaultFilter(bson.D{{"_id", user.Id}}),
		update, false)
	return
}

func (u *UserModelImpl) FindOne(ctx context.Context, id primitive.ObjectID, projections []string) (user *User, err error) {
	err = base_model.QueryOne(ctx, u.Collection(), projections, base_model.DefaultFilter(bson.D{{"_id",
		id}}), &user)

	if err != nil {
		return nil, err
	}

	return
}

func (u *UserModelImpl) List(ctx context.Context, domain base_model.PageDomain, projections []string,
	filter bson.D) (userList []User, total int64, err error) {
	total, err = base_model.Query(ctx, u.Collection(), projections, filter, domain, &userList)
	return
}
