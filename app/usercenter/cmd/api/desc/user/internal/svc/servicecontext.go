package svc

import (
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
