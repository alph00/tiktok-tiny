package main

import (
	"log"
	"net"

	favorite "github.com/alph00/tiktok-tiny/kitex_gen/favorite/favoriteservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"

	// "github.com/bytedance-youthcamp-jbzx/tiktok/cmd/favorite/service"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	// svr := favorite.NewServer(new(FavoriteServiceImpl))

	registryConfig := viper.Read("consul")
	consuladdr := registryConfig.GetString("host") + ":" + registryConfig.GetString("port")
	registry, err := consul.NewConsulRegister(consuladdr)
	if err != nil {
		panic(err)
	}

	favoriteServiceConfig := viper.Read("service")
	favoriteServiceAddr := favoriteServiceConfig.GetString("favorite.host") + ":" + favoriteServiceConfig.GetString("favorite.port")
	addr, err := net.ResolveTCPAddr("tcp", favoriteServiceAddr)
	serviceName := favoriteServiceConfig.GetString("favorite.name")

	svr := favorite.NewServer(new(FavoriteServiceImpl),
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
