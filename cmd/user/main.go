package main

import (
	"log"
	"net"

	"github.com/alph00/tiktok-tiny/cmd/user/service"
	user "github.com/alph00/tiktok-tiny/kitex_gen/user/userservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	registryConfig := viper.Read("consul")
	consuladdr := registryConfig.GetString("host") + ":" + registryConfig.GetString("port")
	registry, err := consul.NewConsulRegister(consuladdr)
	if err != nil {
		panic(err)
	}

	userServiceConfig := viper.Read("service")
	userServiceAddr := userServiceConfig.GetString("user.host") + ":" + userServiceConfig.GetString("user.port")
	addr, err := net.ResolveTCPAddr("tcp", userServiceAddr)
	serviceName := userServiceConfig.GetString("user.name")

	svr := user.NewServer(new(service.UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry),
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
