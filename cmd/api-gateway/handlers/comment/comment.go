package comment

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	jsonutil "github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	"github.com/alph00/tiktok-tiny/kitex_gen/comment"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func init() {
	fmt.Println("------init--------")
}

// TODO 发表评论的业务流程
func CommentAction(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	userId := int64(v.(*model.User).ID)
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	fmt.Println(err)
	if err != nil {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "video_id 不合法", //res.StatusMsg,
			"comment":     nil,
		})
		logger.Error(err)
		return
	}
	//TODO 评论行为判断
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "action_type 不合法",
			"comment":     nil,
		})
		return
	}
	req := new(comment.CommentActionRequest)
	req.Token = token
	req.VideoId = vid
	req.ActionType = int32(actionType)
	req.Id = userId
	//TODO 评论操作 1
	if actionType == 1 {
		commentText := c.Query("comment_text")
		if commentText == "" {
			jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
				"status_code": -1,
				"status_msg":  "comment_text 不能为空",
				"comment":     nil,
			})
			return
		}
		req.CommentText = commentText
	} else if actionType == 2 {
		commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
				"status_code": -1,
				"status_msg":  "comment_id 不合法",
				"comment":     nil,
			})
			return
		}
		req.CommentId = commentID
	}

	res, err := rpc.CommentAction(ctx, req)
	fmt.Println(res)
	fmt.Println(err)

	if res.StatusCode == -1 {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg,
			"comment":     nil,
		})
		return
	}

	jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg,
		"comment":     res.Comment,
	})
}

// TODO 获取评论列表的业务流程
func CommentList(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	userId := int64(v.(*model.User).ID)
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code":  -1,
			"status_msg":   "video_id 不合法",
			"comment_list": nil,
		})
		return
	}
	req := &comment.CommentListRequest{
		Token:   token,
		VideoId: vid,
		Id:      userId,
	}
	res, err := rpc.CommentList(ctx, req)
	if err != nil {
		logger.Error(err)
	}
	if res.StatusCode == -1 {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code":  -1,
			"status_msg":   res.StatusMsg,
			"comment_list": nil,
		})
		return
	}

	jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code":  0,
		"status_msg":   res.StatusMsg,
		"comment_list": res.CommentList,
	})
}
