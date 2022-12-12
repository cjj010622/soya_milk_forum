package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"soya_milk_forum/app/usercenter/cmd/api/internal/config"
	"soya_milk_forum/app/usercenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	Config        config.Config
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
