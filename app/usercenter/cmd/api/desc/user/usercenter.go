package main

import (
	"flag"
	"fmt"

	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/config"
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/handler"
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
