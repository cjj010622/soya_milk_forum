package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MongoDb MongoDb
	JwtAuth JwtAuth
}

type JwtAuth struct {
	AccessSecret string
	AccessExpire int64
}

type MongoDb struct {
	Uri         string
	Db          string
	MaxPoolSize uint64
}
