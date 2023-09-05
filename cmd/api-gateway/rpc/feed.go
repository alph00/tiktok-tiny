package rpc

import (
	"context"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/feed"
	"github.com/alph00/tiktok-tiny/kitex_gen/feed/feedservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	feedClient feedservice.Client
)

func init() {
	resolverConfig := viper.Read("consul")
	consuladdr := resolverConfig.GetString("host") + ":" + resolverConfig.GetString("port")
	resolver, err := consul.NewConsulResolver(consuladdr)
	if err != nil {
		panic(err)
	}

	feedServiceConfig := viper.Read("service")
	serviceName := feedServiceConfig.GetString("feed.name")

	c, err := feedservice.NewClient(
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
	feedClient = c
}

func Feed(ctx context.Context, req *feed.FeedRequest) (*feed.FeedResponse, error) {

	return feedClient.Feed(ctx, req)
}
