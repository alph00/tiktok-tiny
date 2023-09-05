package rpc

import (
	"context"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish/publishservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	publishClient publishservice.Client
)

func init() {
	resolverConfig := viper.Read("consul")
	consuladdr := resolverConfig.GetString("host") + ":" + resolverConfig.GetString("port")
	resolver, err := consul.NewConsulResolver(consuladdr)
	if err != nil {
		panic(err)
	}

	publishServiceConfig := viper.Read("service")
	serviceName := publishServiceConfig.GetString("publish.name")

	c, err := publishservice.NewClient(
		serviceName,
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
	publishClient = c
}

func PublishAction(ctx context.Context, req *publish.PublishActionRequest) (*publish.PublishActionResponse, error) {
	return publishClient.PublishAction(ctx, req)
}

func PublishList(ctx context.Context, req *publish.PublishListRequest) (*publish.PublishListResponse, error) {
	return publishClient.PublishList(ctx, req)
}
