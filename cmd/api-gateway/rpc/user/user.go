package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/kitex_gen/user/userservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/middleware"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	userClient userservice.Client
)

func init() {
	userConfig := viper.Init("user")
	InitUser(&userConfig)
}

// 初始化User服务 etcd客户端
func InitUser(config *viper.Config) {
	fmt.Println("viper初始化user")
	// var config = viper.Init("user")
	fmt.Println("调用了api/rpc/user下的user.go InitUser函数 初始化User服务 客户端")
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		//client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r), // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func Register(ctx context.Context, req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {

	return userClient.Register(ctx, req)
}
func Login(ctx context.Context, req *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	return userClient.Login(ctx, req)
}

func UserInfo(ctx context.Context, req *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	return userClient.UserInfo(ctx, req)
}
