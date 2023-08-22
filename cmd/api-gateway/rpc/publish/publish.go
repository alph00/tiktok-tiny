package rpc

import (
	"context"
	"fmt"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish/publishservice"
	"github.com/alph00/tiktok-tiny/pkg/etcd"
	"github.com/alph00/tiktok-tiny/pkg/middleware"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var (
	publishClient publishservice.Client
)

func init() {
	// video rpc
	videoConfig := viper.Init("video")
	InitVideo(&videoConfig)
}

func InitVideo(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}
	//TODO etcd 客户端相关注册以及配置项
	c := publishservice.MustNewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		//client.WithMuxConnection(1),                        // mux
		//client.WithRPCTimeout(300*time.Second),             // rpc timeout
		//client.WithConnectTimeout(300000*time.Millisecond), // conn timeout
		//client.WithFailureRetry(retry.NewFailurePolicy()),  // retry
		////client.WithSuite(tracing.NewClientSuite()),         // tracer
		client.WithResolver(r), // resolver
		//// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	publishClient = c

}

//func Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
//
//	return publishClient.Feed(ctx, req)
//}

func PublishAction(ctx context.Context, req *publish.PublishActionRequest) (*publish.PublishActionResponse, error) {
	return publishClient.PublishAction(ctx, req)
}

func PublishList(ctx context.Context, req *publish.PublishListRequest) (*publish.PublishListResponse, error) {
	return publishClient.PublishList(ctx, req)
}
