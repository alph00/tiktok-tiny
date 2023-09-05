package rpc

import (
	"context"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/comment"
	"github.com/alph00/tiktok-tiny/kitex_gen/comment/commentservice"
	"github.com/alph00/tiktok-tiny/pkg/consul"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
)

var (
	commentClient commentservice.Client
)

func init() {
	resolverConfig := viper.Read("consul")
	consuladdr := resolverConfig.GetString("host") + ":" + resolverConfig.GetString("port")
	resolver, err := consul.NewConsulResolver(consuladdr)
	if err != nil {
		panic(err)
	}

	commentServiceConfig := viper.Read("service")
	serviceName := commentServiceConfig.GetString("comment.name")

	c, err := commentservice.NewClient(
		serviceName,
		client.WithMuxConnection(1),
		client.WithRPCTimeout(30*time.Second),
		client.WithConnectTimeout(30000*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(resolver), // resolver
		//client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
	commentClient = c
}

func CommentAction(ctx context.Context, req *comment.CommentActionRequest) (*comment.CommentActionResponse, error) {
	return commentClient.CommentAction(ctx, req)
}

func CommentList(ctx context.Context, req *comment.CommentListRequest) (*comment.CommentListResponse, error) {
	return commentClient.CommentList(ctx, req)
}
