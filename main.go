package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"go-im-user-server/rpc/internal/config"
	"go-im-user-server/rpc/internal/server"
	"go-im-user-server/rpc/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/userserver.yaml", "the config file")
var gatewayconfigFile = flag.String("geteway", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	wg := new(sync.WaitGroup)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_server.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		//tttodo http auth
		//	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		reflection.Register(grpcServer)
		//	}
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s.Stop()

		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		s.Start()
	}()

	wg.Add(1)

	var getewayConf gateway.GatewayConf
	conf.MustLoad(*gatewayconfigFile, &getewayConf)
	gw := gateway.MustNewServer(getewayConf)

	go func() {
		defer wg.Done()
		defer gw.Stop()

		time.Sleep(3 * time.Second)
		//	fmt.Printf("Starting rpc server at %s...\n", )
		gw.Start()
	}()

	wg.Wait()
}
