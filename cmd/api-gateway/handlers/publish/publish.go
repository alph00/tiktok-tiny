package publish

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	jsonutil "github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func init() {
	fmt.Println("先执行init函数")
}

// TODO 获取发布列表的操作
func PublishList(ctx context.Context, c *app.RequestContext) {
	token := c.GetString("token")
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id 不合法",
		})
		return
	}
	req := &publish.PublishListRequest{
		Token:  token,
		UserId: uid,
	}
	// fmt.Println(req)
	res, _ := rpc.PublishList(ctx, req)
	// fmt.Println(res)
	if res.StatusCode == -1 {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg,
		})
		return
	}
	jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  res.VideoList,
	})
}

// TODO 上传视频操作
func PublishAction(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	myId := int64(v.(*model.User).ID)
	token := c.PostForm("token")
	title := c.PostForm("title")
	if title == "" {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "标题不能为空",
		})
		return
	}
	// 视频数据
	file, err := c.FormFile("data")
	if err != nil {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "upload error",
		})
		return
	}
	src, _ := file.Open()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, src)
	if err != nil {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "视频上传失败",
		})
		return
	}
	req := &publish.PublishActionRequest{
		Token: token,
		Title: title,
		Data:  buf.Bytes(),
		Id:    myId,
	}
	res, _ := rpc.PublishAction(ctx, req)
	if res.StatusCode == -1 {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg,
		})
		return
	}
	jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg,
	})

}
