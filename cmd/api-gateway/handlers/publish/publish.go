package handler

import (
	"bytes"
	"context"
	"fmt"
	rpc "github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc/publish"
	"github.com/alph00/tiktok-tiny/internal/response"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"io"
	"net/http"
	"strconv"
)

//func Feed(ctx context.Context, c *app.RequestContext) {
//	token := c.Query("token")
//	latestTime := c.Query("latest_time")
//	var timestamp int64 = 0
//	if latestTime != "" {
//		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
//	} else {
//		timestamp = time.Now().UnixMilli()
//	}
//
//	req := &kitex.FeedRequest{
//		LatestTime: timestamp,
//		Token:      token,
//	}
//	res, _ := rpc.Feed(ctx, req)
//	if res.StatusCode == -1 {
//		c.JSON(http.StatusOK, response.Feed{
//			Base: response.Base{
//				StatusCode: -1,
//				StatusMsg:  res.StatusMsg,
//			},
//		})
//		return
//	}
//	c.JSON(http.StatusOK, response.Feed{
//		Base: response.Base{
//			StatusCode: 0,
//			StatusMsg:  res.StatusMsg,
//		},
//		VideoList: res.VideoList,
//	})
//}

// TODO 获取发布列表的操作
func PublishList(ctx context.Context, c *app.RequestContext) {
	token := c.GetString("token")

	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}
	req := &publish.PublishListRequest{
		Token:  token,
		UserId: uid,
	}
	fmt.Println(req)
	res, _ := rpc.PublishList(ctx, req)
	fmt.Println(res)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.PublishList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: res.VideoList,
	})
}

// TODO 上传视频操作
func PublishAction(ctx context.Context, c *app.RequestContext) {
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户鉴权失败，token为空",
			},
		})
		return
	}
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "标题不能为空",
			},
		})
		return
	}
	// 视频数据
	file, err := c.FormFile("data")
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "上传视频加载失败",
			},
		})
		return
	}
	src, err := file.Open()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "视频上传失败",
			},
		})
		return
	}

	req := &publish.PublishActionRequest{
		Token: token,
		Title: title,
		Data:  buf.Bytes(),
	}
	res, _ := rpc.PublishAction(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.PublishAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
	})
}
