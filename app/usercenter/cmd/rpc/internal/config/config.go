package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MongoDb MongoDb
}

type MongoDb struct {
	Uri         string
	Db          string
	MaxPoolSize uint64
}
