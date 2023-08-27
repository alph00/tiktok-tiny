package rpc

import (
	"context"
	"fmt"
	"github.com/alph00/tiktok-tiny/kitex_gen/comment"
	"github.com/alph00/tiktok-tiny/kitex_gen/comment/commentservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/middleware"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	commentClient commentservice.Client
)

func init() {
	// comment rpc
	commentConfig := viper.Init("comment")
	InitComment(&commentConfig)
}

func InitComment(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	//TODO 客户端
	c, err := commentservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		//client.WithMuxConnection(1),                       // mux
		//client.WithRPCTimeout(30*time.Second),             // rpc timeout
		//client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		//client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		//client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r), // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
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
