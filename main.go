package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/heyehang/go-im-pkg/tlog"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"go-im-user-server/internal/config"
	"go-im-user-server/internal/pyroscope"
	"go-im-user-server/internal/server"
	"go-im-user-server/internal/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configFile        = flag.String("f", "etc/userserver.yaml", "the config file")
	gatewayconfigFile = flag.String("gateway", "etc/gateway.yaml", "the config file")
)

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	fileWriter := logx.Reset()
	writer, err := tlog.NewMultiWriter(fileWriter)
	logx.Must(err)
	logx.SetWriter(writer)
	ctx := svc.NewServiceContext(c)
	wg := new(sync.WaitGroup)
	waitGrpcStart := make(chan struct{}, 1)
	gs := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_server.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		//tttodo http auth
		//	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		reflection.Register(grpcServer)
		//	}x
		waitGrpcStart <- struct{}{}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer gs.Stop()
		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		gs.Start()
	}()
	<-waitGrpcStart
	wg.Add(1)
	var getewayConf gateway.GatewayConf
	conf.MustLoad(*gatewayconfigFile, &getewayConf)
	gw := gateway.MustNewServer(getewayConf)
	go func() {
		defer wg.Done()
		defer gw.Stop()
		time.Sleep(3 * time.Second)
		fmt.Printf("Starting rpc server at %s...\n", getewayConf.Host+fmt.Sprintf("%d", getewayConf.Prometheus.Port))
		logx.Slowf("Starting rpc server at ...\n", logx.Field("addr", getewayConf.Host+fmt.Sprintf("%d", getewayConf.Prometheus.Port)))
		gw.Start()
	}()
	cancenCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pyroscope.Start(cancenCtx, wg, c.Name, c.PyroscopeAddr, nil, true)
	defer pyroscope.Closed()
	wg.Wait()
}
