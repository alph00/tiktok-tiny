package rpc

import (
	"context"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/relation"
	"github.com/alph00/tiktok-tiny/kitex_gen/relation/relationservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	relationClient relationservice.Client
)

func init() {
	resolverConfig := viper.Read("consul")
	consuladdr := resolverConfig.GetString("host") + ":" + resolverConfig.GetString("port")
	resolver, err := consul.NewConsulResolver(consuladdr)
	if err != nil {
		panic(err)
	}

	relationServiceConfig := viper.Read("service")
	serviceName := relationServiceConfig.GetString("relation.name")

	c, err := relationservice.NewClient(
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
	relationClient = c
}
func RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	return relationClient.RelationAction(ctx, req)
}

// FollowList implements the RelationServiceImpl interface.
func FollowList(ctx context.Context, req *relation.RelationFollowListRequest) (resp *relation.RelationFollowListResponse, err error) {
	return relationClient.FollowList(ctx, req)
}

// FollowerList implements the RelationServiceImpl interface.
func FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	return relationClient.FollowerList(ctx, req)
}

// FriendList implements the RelationServiceImpl interface.
func FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	return relationClient.FriendList(ctx, req)
}
