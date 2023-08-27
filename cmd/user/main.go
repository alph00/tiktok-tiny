package main

import (
	"fmt"
	"net"

	"github.com/alph00/tiktok-tiny/cmd/user/service"
	"github.com/alph00/tiktok-tiny/kitex_gen/user/userservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/middleware"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/alph00/tiktok-tiny/pkg/zap"

	// "github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

var (
	config      = viper.Init("user")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
	logger      = zap.InitLogger()
)

func init() {
	service.Init(signingKey)
}
func main() {
	// svr := user.NewServer(new(UserServiceImpl))

	// err := svr.Run()

	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
	} else {
		// TODO:这里打印服务建立成功
		fmt.Println("etcd服务建立成功")
	}

	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
	} else {
		fmt.Println("网络连接建立成功")
	}
	fmt.Printf("server: %s:%d\n", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	fmt.Printf("etcd: %s:%d\n", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	// 初始化etcd
	s := userservice.NewServer(new(service.UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithRegistry(r),
		//server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err := s.Run(); err != nil {
		logger.Fatalf("%v stopped with error: %v", serviceName, err.Error())
	} else {
		fmt.Println("user微服务服务端建立成功")
	}

}
