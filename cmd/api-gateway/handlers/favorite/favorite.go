package favorite

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	kitex "github.com/alph00/tiktok-tiny/kitex_gen/favorite"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	fmt.Println("调用了api/handler/favorite下的user.go FavoriteAction函数")

	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	// TODO:这里应该不需要判断是否合法吧
	id, err := strconv.ParseInt(videoId, 10, 64)

	if err != nil {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "video_id 不合法", //res.StatusMsg,

		})
		return
	}
	acttype, err := strconv.Atoi(actionType)
	if err != nil {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "actiontype 不合法", //res.StatusMsg,

		})
		return
	}
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	fmt.Printf("登录用户ID为：%d\n", int64(v.(*model.User).ID))
	req := &kitex.FavoriteActionRequest{
		UserId:     int64(v.(*model.User).ID),
		Token:      token,
		VideoId:    id,
		ActionType: int32(acttype),
	}
	fmt.Println("进行rpc注册，调用handlers下的FavoriteAction函数")
	res, _ := rpc.FavoriteAction(ctx, req)
	fmt.Printf("进行rpc注册，调用handlers下的FavoriteAction函数完成,%v\n", res)
	if res.StatusCode == -1 {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg, //res.StatusMsg,
		})
		return
	}
	myutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg, //res.StatusMsg,
	})

}
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	fmt.Println("调用了api/handler/favorite下的user.go FavoriteList函数")

	token := c.Query("token")
	userId := c.Query("user_id")
	// actionType := c.Query("action_type")
	// TODO:这里应该不需要判断是否合法吧
	UserId, err := strconv.ParseInt(userId, 10, 64)

	if err != nil {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id 不合法", //res.StatusMsg,

		})
		return
	}
	req := &kitex.FavoriteListRequest{
		UserId: UserId,
		Token:  token,
	}
	fmt.Println("进行rpc注册，调用rpc/user下的FavoriteList函数")
	res, _ := rpc.FavoriteList(ctx, req)
	fmt.Printf("进行rpc注册，调用rpc/user下的FavoriteList函数完成,%v\n", res)
	if res.StatusCode == -1 {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg, //res.StatusMsg,
			"video_list":  nil,
		})
		return
	}
	myutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg, //res.StatusMsg,
		"video_list":  res.VideoList,
	})

}
