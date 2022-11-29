package base_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"soya_milk_forum/common/core"
	"time"
)

type BaseBean struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	IsDeleted int8               `bson:"is_deleted"`
	CreatedAt int64              `bson:"created_at,omitempty"`
	CreatedBy int64              `bson:"created_by,omitempty"`
	UpdatedAt int64              `bson:"updated_at,omitempty"`
	UpdatedBy int64              `bson:"updated_by,omitempty"`
}

type PageDomain struct {
	PageSize      int64  `json:"page_size,default=20,range=[0:150]"`
	PageNum       int64  `json:"page_num,default=1"`
	OrderByColumn string `json:"order_by_column,optional"`
	Sort          string `json:"sort,default=asc,options=asc|desc"`
}

type Interface interface {
	TableName() string
	Collection() *mongo.Collection
	DelMany(ctx context.Context, ids []primitive.ObjectID, deletedBy int64) error
}

func DefaultFilter(filter bson.D) bson.D {
	filter = append(filter, bson.E{"is_deleted", 0})
	return filter
}

// InitTime 如果CreatedAt与UpdatedAt为默认零值则赋值当前时间
func InitTime(t any) {
	tVal := reflect.ValueOf(t)
	tVal = core.ElemValueIfPointer(tVal)
	cAt, uAt := "CreatedAt", "UpdatedAt"
	now := reflect.ValueOf(time.Now().Unix())

	if core.IsDefaultZero(tVal.FieldByName(cAt).Interface()) {
		tVal.FieldByName(cAt).Set(now)
	}
	if core.IsDefaultZero(tVal.FieldByName(uAt).Interface()) {
		tVal.FieldByName(uAt).Set(now)
	}
}

func DefaultDelMany(ctx context.Context, collection *mongo.Collection, fieldName string,
	ids []primitive.ObjectID, deletedBy int64) (err error) {

	if len(ids) == 0 {
		return
	}

	fa := bson.A{}
	for _, id := range ids {
		fa = append(fa, bson.D{{fieldName, id}})
	}
	filter := bson.D{{"$or", fa}}

	update := bson.D{{"$set", bson.D{{"is_deleted", 1},
		{"updated_by", deletedBy},
		{"updated_at", time.Now().Unix()}}}}
	_, err = collection.UpdateMany(ctx, filter, update)

	return
}

// Update 默认初始值将不会更新
func Update(ctx context.Context, collection *mongo.Collection, filter interface{}, update any) (*mongo.UpdateResult, error) {

	return collection.UpdateOne(ctx,
		filter,
		bson.D{{"$set", update}})
}

// UpsertOne 如果upsert=true则有则修改，无则新增。
func UpsertOne(ctx context.Context, collection *mongo.Collection, filter bson.D,
	update any, upsert bool) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(upsert)
	up := bson.D{{"$set", update}}

	result, err := collection.UpdateOne(ctx, filter, up, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// QueryOne 简化查询集合单条信息函数
func QueryOne(ctx context.Context, collection *mongo.Collection, projections []string, filter bson.D,
	res interface{}) error {

	proBson := bson.M{}
	for _, v := range projections {
		proBson[v] = 1
	}
	opts := options.FindOneOptions{Projection: proBson}
	return collection.FindOne(ctx, filter, &opts).Decode(res)
}

// Query 查询多个参数
func Query(ctx context.Context, collection *mongo.Collection, projections []string, filter bson.D,
	domain PageDomain, res any) (total int64, err error) {
	proBson := bson.M{}
	for _, v := range projections {
		proBson[v] = 1
	}
	opts := options.FindOptions{Projection: proBson}
	// 处理分页
	handlerPageDomain(&opts, domain)

	cur, err := collection.Find(ctx, filter, &opts)
	if err != nil {
		return 0, err
	}

	total, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return total, cur.All(ctx, res)
}

// 处理分页参数
func handlerPageDomain(opts *options.FindOptions, domain PageDomain) {
	// 默认每页大小
	opts.SetLimit(20)
	// 处理排序
	if domain.OrderByColumn != "" {
		if domain.Sort != "" && domain.Sort != "asc" {
			opts.SetSort(bson.D{{domain.OrderByColumn, -1}})
		} else {
			opts.SetSort(bson.D{{domain.OrderByColumn, 1}})
		}
	}
	// 处理分页
	if domain.PageSize != 0 {
		opts.SetLimit(domain.PageSize)
		if domain.PageNum > 0 {
			opts.SetSkip((domain.PageNum - 1) * domain.PageSize)
		}
	}
}
