package base_model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"soya_milk_forum/common/core"
)

// GetFieldNamesByBsonTag 获取给定obj的字段Tag的bson值，如果allField为true则获取所有
func GetFieldNamesByBsonTag(obj any, allField bool) (res []string) {
	objV := reflect.ValueOf(obj)
	objV = core.ElemValueIfPointer(objV)
	objT := objV.Type()
	res = append(res, "_id")
	for i := 0; i < objV.NumField(); i++ {
		if allField || objV.Field(i).Interface().(int8) == 1 {
			res = append(res, objT.Field(i).Tag.Get("bson"))
		}
	}

	return res
}

func ToLookupStage(from, localField, foreignField, as string) bson.D {
	return bson.D{{"$lookup", bson.D{
		{"from", from},
		{"localField", localField},
		{"foreignField", foreignField},
		{"as", as},
	}}}
}

func ToUnwindStage(pv string, flag bool) bson.D {
	return bson.D{{"$unwind", bson.D{
		{"path", pv},
		{"preserveNullAndEmptyArrays", flag},
	}}}
}

func ToProject(project bson.M, v int, projections []string, prefix string) bson.M {
	for _, p := range projections {
		project[prefix+p] = v
	}
	return project
}

// HandlePageDomain 处理分页参数
func HandlePageDomain(domain PageDomain) []bson.D {
	var pageStage []bson.D
	// 默认每页大小
	limitStage := bson.D{{"$limit", 20}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	skipStage := bson.D{{"$skip", 0}}

	// 处理排序
	if domain.OrderByColumn != "" {
		if domain.Sort != "" && domain.Sort != "asc" {
			sortStage = bson.D{{"$sort", bson.D{{domain.OrderByColumn, -1}}}}
		} else {
			sortStage = bson.D{{"$sort", bson.D{{domain.OrderByColumn, 1}}}}
		}
	}
	// 处理分页
	if domain.PageSize != 0 {
		limitStage = bson.D{{"$limit", domain.PageSize}}
		if domain.PageNum > 0 {
			skipStage = bson.D{{"$skip", (domain.PageNum - 1) * domain.PageSize}}
		}
	}

	pageStage = append(pageStage, sortStage, skipStage, limitStage)

	return pageStage
}

// ExecutePipeline 执行分页聚合查询
func ExecutePipeline(ctx context.Context, collection *mongo.Collection, pipeline mongo.Pipeline,
	res any, domain PageDomain) (total int64, err error) {

	cur, err := collection.Aggregate(ctx, append(pipeline, HandlePageDomain(domain)...))
	if err != nil {
		return
	}

	err = cur.All(ctx, res)
	if err != nil {
		return
	}

	total, err = PipelineTotal(ctx, collection, pipeline)

	return
}

func PipelineTotal(ctx context.Context, collection *mongo.Collection,
	pipeline mongo.Pipeline) (total int64, err error) {
	// 查询总数
	countStage := bson.D{{"$count", "total_documents"}}
	curTotal, err := collection.Aggregate(ctx, append(pipeline, countStage))
	if err != nil {
		return
	}

	var counts []countTotal
	err = curTotal.All(ctx, &counts)
	if err != nil {
		return
	}

	if len(counts) > 0 {
		total = counts[0].TotalDocuments
	}
	return
}

type countTotal struct {
	TotalDocuments int64 `bson:"total_documents"`
}

func idsToBsonA(ids []primitive.ObjectID, field string) (res bson.A) {
	for _, id := range ids {
		res = append(res, bson.D{{field, id}})
	}
	return
}
