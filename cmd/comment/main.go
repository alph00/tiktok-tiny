package main

import (
	"fmt"
	comment "github.com/alph00/tiktok-tiny/kitex_gen/comment/commentservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"net"
)

var (
	config      = viper.Init("comment")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
)

func init() {
	Init(signingKey)
}

func main() {
	////TODO etcd服务注册
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
	s := comment.NewServer(new(CommentServiceImpl),
		//server.WithServiceAddr(addr),
		//server.WithMiddleware(middleware.CommonMiddleware),
		//server.WithMiddleware(middleware.ServerMiddleware),
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

	//r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	////创建一个etcd的注册，ectd占127.0.0.1的2379端口
	//if err != nil {
	//	fmt.Println(err)
	//}
	//svr := comment.NewServer(new(CommentServiceImpl),
	//	server.WithRegistry(r), //将server注册到etcd中
	//	server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
	//		ServiceName: serviceName,
	//	}),
	//)
	//
	//err = svr.Run()
	//
	//if err != nil {
	//	log.Println(err.Error())
	//}

	//for {
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	//	resp, err := c.CommentAction(ctx, &req.CommentActionRequest{
	//		Token:       "1111111",
	//		VideoId:     1,
	//		ActionType:  1,
	//		CommentText: "test",
	//	})
	//	cancel()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	time.Sleep(time.Second)
	//	fmt.Println(resp)
	//}

}
