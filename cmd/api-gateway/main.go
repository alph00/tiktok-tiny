package main

import (
	"context"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/message"
	mw "github.com/alph00/tiktok-tiny/pkg/mw/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

var (
	apiServerAddr = "0.0.0.0:8080"
	// consulAddress = "0.0.0.0:2379"
)

func main() {
	r := server.New(
		server.WithHostPorts(apiServerAddr),
		server.WithHandleMethodNotAllowed(true),
	)

	r.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "no route")
	})
	// r.NoMethod(func(ctx context.Context, c *app.RequestContext) {
	// 	c.String(consts.StatusOK, "no method")
	// })

	// the jwt middleware
	mw.InitJwt()

	// router group for jwt authentication
	var tiktok_tiny, usr, fed, msg, rel, pub, fav, comm, commAction, commList *route.RouterGroup
	tiktok_tiny = r.Group("/douyin")
	usr = tiktok_tiny.Group("/user")
	{
		groupNoAuthUse(usr)
		// usr.GET("/", user.UserInfo)
		// usr.POST("/register/", user.Register)
		usr.POST("/login/", mw.JwtMiddleware.LoginHandler)
	}
	fed = tiktok_tiny.Group("/feed")
	{
		groupNoAuthUse(fed)
		// fed.GET("/", feed.Feed)
	}
	msg = tiktok_tiny.Group("/message")
	{
		groupAuthUse(msg)
		msg.GET("/chat/", message.MessageChat)
		msg.POST("/action/", message.MessageAction)
	}
	rel = tiktok_tiny.Group("/relation")
	{
		groupAuthUse(rel)
		// rel.GET("/follower/list/", relation.FollowerList)
		// rel.GET("/follow/list/", relation.FollowList)
		// rel.GET("/friend/list/", relation.FriendList)
		// rel.POST("/action/", relation.RelationAction)
	}
	pub = tiktok_tiny.Group("/publish")
	{
		groupAuthUse(pub)
		// pub.GET("/list/", publish.PublishList)
		// pub.POST("/action/", publish.PublishAction)
	}
	fav = tiktok_tiny.Group("/favorite")
	{
		groupAuthUse(fav)
		// fav.POST("/action/", favorite.FavoriteAction)
		// fav.GET("/list/", favorite.FavoriteList)
	}
	comm = tiktok_tiny.Group("/comment")
	commAction = comm.Group("/action")
	{
		groupAuthUse(commAction)
		// commAction.POST("/", comment.CommentAction)
	}
	commList = comm.Group("/list")
	{
		groupNoAuthUse(commList)
		// commList.GET("/", comment.CommentList)
	}
	r.Spin()
}

func groupNoAuthUse(group ...*route.RouterGroup) {
	for _, g := range group {
		g.Use(recovery.Recovery())
	}
}
func groupAuthUse(group ...*route.RouterGroup) {
	for _, g := range group {
		g.Use(recovery.Recovery(), mw.JwtMiddleware.MiddlewareFunc())
	}
}
