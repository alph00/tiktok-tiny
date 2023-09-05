package main

import (
	"fmt"
	"log"
	"net"

	"github.com/alph00/tiktok-tiny/cmd/feed/service"
	feed "github.com/alph00/tiktok-tiny/kitex_gen/feed/feedservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	//TODO 读取feed的配置，服务注册到Consul中
	registryConfig := viper.Read("consul")
	consuladdr := registryConfig.GetString("host") + ":" + registryConfig.GetString("port")
	registry, err := consul.NewConsulRegister(consuladdr)
	if err != nil {
		panic(err)
	}

	feedServiceConfig := viper.Read("service")
	feedServiceAddr := feedServiceConfig.GetString("feed.host") + ":" + feedServiceConfig.GetString("feed.port")
	addr, err := net.ResolveTCPAddr("tcp", feedServiceAddr)
	serviceName := feedServiceConfig.GetString("feed.name")
	fmt.Printf("serviceName: %v\n", serviceName)
	fmt.Printf("consuladdr: %v\n", consuladdr)

	svr := feed.NewServer(new(service.FeedServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry),
		//TODO 限流开启不???
		//server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		// server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
