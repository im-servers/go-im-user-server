package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"go-im-user-server/internal/config"
	"go-im-user-server/internal/server"
	"go-im-user-server/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/heyehang/go-im-pkg/tlog"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var profile *pyroscope.Profiler
var configFile = flag.String("f", "etc/userserver.yaml", "the config file")
var gatewayconfigFile = flag.String("gateway", "etc/gateway.yaml", "the config file")

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
	gs := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_server.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		//tttodo http auth
		//	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		reflection.Register(grpcServer)
		//	}
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer gs.Stop()

		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		logx.Slowf("Starting rpc server at ...\n", logx.Field("addr", c.ListenOn))
		gs.Start()
	}()

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if profile != nil {
				_ = profile.Stop()
			}
		}()
		startPyroscope()
	}()

	wg.Wait()
}

func startPyroscope() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	var err error
	profile, err = pyroscope.Start(pyroscope.Config{
		ApplicationName: "go-im-user-server",
		// replace this with the address of pyroscope server
		ServerAddress: "http://172.16.0.15:4040",
		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,
		// optionally, if authentication is enabled, specify the API key:
		// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),
		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		SampleRate: 200,
	})
	if err != nil {
		panic(err)
		return
	}
}
