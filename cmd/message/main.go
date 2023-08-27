package main

import (
	"fmt"
	"log"
	"net"

	"github.com/alph00/tiktok-tiny/cmd/message/service"
	message "github.com/alph00/tiktok-tiny/kitex_gen/message/messageservice"
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

	messageServiceConfig := viper.Read("service")
	messageServiceAddr := messageServiceConfig.GetString("message.host") + ":" + messageServiceConfig.GetString("message.port")
	addr, err := net.ResolveTCPAddr("tcp", messageServiceAddr)
	serviceName := messageServiceConfig.GetString("message.name")
	fmt.Printf("serviceName: %v\n", serviceName)
	fmt.Printf("consuladdr: %v\n", consuladdr)

	svr := message.NewServer(new(service.MessageServiceImpl),
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
