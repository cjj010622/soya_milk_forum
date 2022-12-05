package stores

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"reflect"
	"soya_milk_forum/common/core"
	"time"
)

var KeysLenError = errors.New("keys长度错误")

type (
	// BatchGetCallback 缓存中不存在 key 则会调用此函数
	BatchGetCallback func(ctx context.Context, key string) (res any, err error)

	RedisModel interface {
		// BatchSetEX 批量保存数据并设置过期时间，过期时间单位为秒，len(keys)!=len(arr) 则会抛出异常
		BatchSetEX(ctx context.Context, keys []string, arr any, expire int64) error
		// BatchGet 批量获取缓存中的数据并反序列化到res中，如果 key 不存在则会调用回调函数获取数据
		BatchGet(ctx context.Context, keys []string, res interface{}, callback BatchGetCallback) error
		// BatchDel 批量删除
		BatchDel(ctx context.Context, keys []string) error
	}

	RedisModelImpl struct {
		client *redis.Client
	}
)

func (r *RedisModelImpl) BatchDel(ctx context.Context, keys []string) error {
	pipeline := r.client.Pipeline()
	for i := 0; i < len(keys); i++ {
		pipeline.Del(ctx, keys[i])
	}
	_, err := pipeline.Exec(ctx)

	return err
}

func (r *RedisModelImpl) BatchGet(ctx context.Context, keys []string, res interface{}, callback BatchGetCallback) error {
	// 检查 res 类型是否为 切片指针
	if err := checkRes(res); err != nil {
		return err
	}

	pipeline := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i := 0; i < len(keys); i++ {
		cmds[i] = pipeline.Get(ctx, keys[i])
	}
	_, err := pipeline.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	// 解析结果
	resVal := reflect.ValueOf(res)
	sliceVal := resVal.Elem()

	elementType := sliceVal.Type().Elem()
	var index int

	for i, v := range cmds {
		val, err := v.Result()
		// key 不存在则使用回调函数获取
		if err == redis.Nil {
			val, err = r.callbackAndSet(ctx, keys[i], callback)
		}
		if err != nil {
			return err
		}

		// 需要扩容
		if sliceVal.Len() == index {
			newElem := reflect.New(elementType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
		}

		curElem := sliceVal.Index(index).Addr().Interface()
		err = json.Unmarshal([]byte(val), &curElem)
		if err != nil {
			return err
		}
		index++
	}

	resVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}

func (r *RedisModelImpl) callbackAndSet(ctx context.Context, key string,
	callback BatchGetCallback) (val string, err error) {

	if callback == nil {
		return "", errors.New("key: " + key + " 不存在，同时为设置回调函数")
	}

	cv, err := callback(ctx, key)
	if err != nil {
		return
	}

	bs, err := json.Marshal(cv)
	_, err = r.client.Set(ctx, key, bs, 3600*time.Second).Result()
	if err != nil {
		return
	}

	return string(bs), nil
}

func checkRes(res interface{}) error {
	resVal := reflect.ValueOf(res)
	if resVal.Kind() != reflect.Pointer {
		return errors.New("res 必须为指针")
	}

	resVal = resVal.Elem()
	if resVal.Kind() != reflect.Slice {
		return errors.New("res 必须为指针类型的切片")
	}

	return nil
}

func (r *RedisModelImpl) BatchSetEX(ctx context.Context, keys []string, arr any, expire int64) error {
	if arr == nil {
		return errors.New("arr 不能为空")
	}
	arrVal := reflect.ValueOf(arr)
	arrVal = core.ElemValueIfPointer(arrVal)
	if arrVal.Len() != len(keys) {
		return KeysLenError
	}

	pipeline := r.client.Pipeline()
	for i := 0; i < arrVal.Len(); i++ {
		bs, err := json.Marshal(arrVal.Index(i).Interface())
		if err != nil {
			return err
		}
		pipeline.SetEX(ctx, keys[i], bs, time.Duration(expire)*time.Second)
	}

	_, err := pipeline.Exec(ctx)

	return err
}

func NewRedisModel(client *redis.Client) RedisModel {
	return &RedisModelImpl{
		client: client,
	}
}
