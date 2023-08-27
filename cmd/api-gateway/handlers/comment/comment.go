package comment

import (
	"context"
	"fmt"
	rpc "github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc/comment"
	"github.com/alph00/tiktok-tiny/internal/response"
	"github.com/alph00/tiktok-tiny/kitex_gen/comment"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"
)

// TODO 发表评论的业务流程
func CommentAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
			Comment: nil,
		})
		logger.Error(err)
		return
	}
	//TODO 评论行为判断
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
			Comment: nil,
		})
		return
	}
	req := new(comment.CommentActionRequest)
	req.Token = token
	req.VideoId = vid
	req.ActionType = int32(actionType)
	//TODO 评论操作 1
	if actionType == 1 {
		commentText := c.Query("comment_text")
		if commentText == "" {
			c.JSON(http.StatusOK, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "comment_text 不能为空",
				},
				Comment: nil,
			})
			return
		}
		req.CommentText = commentText
	} else if actionType == 2 {
		commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "comment_id 不合法",
				},
				Comment: nil,
			})
			return
		}
		req.CommentId = commentID
	}
	res, err := rpc.CommentAction(ctx, req)
	fmt.Println(res)
	fmt.Println(err)

	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			Comment: nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.CommentAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		Comment: res.Comment,
	})
}

// TODO 获取评论列表的业务流程
func CommentList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.CommentList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
			CommentList: nil,
		})
		return
	}
	req := &comment.CommentListRequest{
		Token:   token,
		VideoId: vid,
	}
	res, _ := rpc.CommentList(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			CommentList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.CommentList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		CommentList: res.CommentList,
	})
}
