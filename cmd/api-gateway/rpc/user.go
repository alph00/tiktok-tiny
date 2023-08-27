package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/kitex_gen/user/userservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	userClient userservice.Client
)

func init() {
	InitUser()
}

// 初始化User服务 consul客户端
func InitUser() {
	fmt.Println("viper初始化user")
	// var config = viper.Init("user")
	fmt.Println("调用了api/rpc/user下的user.go InitUser函数 初始化User服务 客户端")
	resolverConfig := viper.Read("consul")
	consuladdr := resolverConfig.GetString("host") + ":" + resolverConfig.GetString("port")
	resolver, err := consul.NewConsulResolver(consuladdr)
	if err != nil {
		panic(err)
	}
	userConfig := viper.Read("service")
	serviceName := userConfig.GetString("user.name")

	c, err := userservice.NewClient(
		serviceName,
		// client.WithMiddleware(middleware.CommonMiddleware),
		// client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(30*time.Second),
		client.WithConnectTimeout(30000*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(resolver), // resolver
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
