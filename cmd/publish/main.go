package main

import (
	"fmt"
	publish "github.com/alph00/tiktok-tiny/kitex_gen/publish/publishservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/middleware"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

var (
	config      = viper.Init("video")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
)

func init() {
	Init(signingKey)
}

func main() {

	//TODO etcd服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		logger.Errorf(err.Error())
	}

	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		logger.Errorf(err.Error())
	}
	fmt.Println(addr)
	//TODO 服务端相关的配置项以及注册
	s := publish.NewServer(new(PublishServiceImpl),
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithRegistry(r),
		////server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		//server.WithMuxTransport(),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err := s.Run(); err != nil {
		fmt.Println(err)
		logger.Errorf("%v stopped with error: %v", serviceName, err.Error())
	}

	if err != nil {
		log.Println(err.Error())
	}
}
