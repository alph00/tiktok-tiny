package main

import (
	"crypto/tls"
	"fmt"
	commentHandler "github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/comment"
	publishHandler "github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/publish"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/gzip"
)

var (
	apiConfig     = viper.Init("api")
	apiServerName = apiConfig.Viper.GetString("server.name")
	apiServerAddr = fmt.Sprintf("%s:%d", apiConfig.Viper.GetString("server.host"), apiConfig.Viper.GetInt("server.port"))
	etcdAddress   = fmt.Sprintf("%s:%d", apiConfig.Viper.GetString("Etcd.host"), apiConfig.Viper.GetInt("Etcd.port"))
	signingKey    = apiConfig.Viper.GetString("JWT.signingKey")
	serverTLSKey  = apiConfig.Viper.GetString("Hertz.tls.keyFile")
	serverTLSCert = apiConfig.Viper.GetString("Hertz.tls.certFile")
)

func registerGroup(hz *server.Hertz) {
	douyin := hz.Group("/douyin")
	{

		publishGroup := douyin.Group("/publish")
		{
			publishGroup.GET("/list/", publishHandler.PublishList)
			publishGroup.POST("/action/", publishHandler.PublishAction)
		}
		//douyin.GET("/feed", handler.Feed)

		comment := douyin.Group("/comment")
		{
			comment.POST("/action/", commentHandler.CommentAction)
			comment.GET("/list/", commentHandler.CommentList)
		}
	}
}

func InitHertz() *server.Hertz {

	opts := []config.Option{server.WithHostPorts(apiServerAddr)}

	// 网络库
	hertzNet := standard.NewTransporter
	//if apiConfig.Viper.GetBool("Hertz.useNetPoll") {
	//	hertzNet = netpoll.NewTransporter
	//}
	opts = append(opts, server.WithTransport(hertzNet))

	// TLS & Http2
	// https://github.com/cloudwego/hertz-examples/blob/main/protocol/tls/main.go
	tlsConfig := tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	if apiConfig.Viper.GetBool("Hertz.tls.enable") {
		if len(serverTLSKey) == 0 {
			panic("not found tiktok_tls_key in configuration")
		}
		if len(serverTLSCert) == 0 {
			panic("not found tiktok_tls_cert in configuration")
		}

		cert, err := tls.LoadX509KeyPair(serverTLSCert, serverTLSKey)
		if err != nil {
			logger.Error(err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		opts = append(opts, server.WithTLS(&tlsConfig))

		if alpn := apiConfig.Viper.GetBool("Hertz.tls.ALPN"); alpn {
			opts = append(opts, server.WithALPN(alpn))
		}
	} else if apiConfig.Viper.GetBool("Hertz.http2.enable") {
		opts = append(opts, server.WithH2C(apiConfig.Viper.GetBool("Hertz.http2.enable")))
	}

	hz := server.Default(opts...)

	hz.Use(
		// secure.New(
		// 	secure.WithSSLHost(apiServerAddr),
		// 	secure.WithSSLRedirect(true),
		// ),	// TLS
		//middleware.TokenAuthMiddleware(*jwt.NewJWT([]byte(signingKey)),
		//	"/douyin/user/register/",
		//	"/douyin/user/login/",
		//	"/douyin/feed",
		//	"/douyin/favorite/list/",
		//	"/douyin/publish/list/",
		//	"/douyin/comment/list/",
		//	"/douyin/relation/follower/list/",
		//	"/douyin/relation/follow/list/",
		//), // 用户鉴权中间件
		//middleware.TokenLimitMiddleware(), //限流中间件
		//middleware.AccessLog(),
		gzip.Gzip(gzip.DefaultCompression),
	)
	return hz
}

func main() {
	hz := InitHertz()

	// add handler
	registerGroup(hz)

	hz.Spin()
}
