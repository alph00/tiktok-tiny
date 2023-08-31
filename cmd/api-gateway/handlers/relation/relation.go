package relation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	"github.com/alph00/tiktok-tiny/kitex_gen/relation"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func FriendList(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	myId := int64(v.(*model.User).ID)
	useridString := c.Query("user_id")
	userid, err := strconv.ParseInt(useridString, 10, 64)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id格式错误",
		})
		return
	}
	if myId != userid {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id不匹配",
		})
		return
	}
	response, err := rpc.FriendList(ctx, &relation.RelationFriendListRequest{
		UserId: myId,
		Token:  c.Query("token"),
	})
	fmt.Printf("qfsresponse: %v\n", response)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusInternalServerError), c, utils.H{
			"status_code": -1, //res.StatusCode
			"status_msg":  response.StatusMsg,
		})
		return
	}

	myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
		"status_code": 0, //res.StatusCode
		"status_msg":  response.StatusMsg,
		"user_list":   response.UserList,
	})
}
func FollowerList(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	myId := int64(v.(*model.User).ID)
	useridString := c.Query("user_id")
	userid, err := strconv.ParseInt(useridString, 10, 64)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id格式错误",
		})
		return
	}
	if myId != userid {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id不匹配",
		})
		return
	}
	response, err := rpc.FollowerList(ctx, &relation.RelationFollowerListRequest{
		UserId: myId,
		Token:  c.Query("token"),
	})
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusInternalServerError), c, utils.H{
			"status_code": -1, //res.StatusCode
			"status_msg":  response.StatusMsg,
		})
		return
	}

	myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
		"status_code": 0, //res.StatusCode
		"status_msg":  response.StatusMsg,
		"user_list":   response.UserList,
	})
}
func FollowList(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	myId := int64(v.(*model.User).ID)
	useridString := c.Query("user_id")
	userid, err := strconv.ParseInt(useridString, 10, 64)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id格式错误",
		})
		return
	}
	if myId != userid {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id不匹配",
		})
		return
	}
	response, err := rpc.FollowList(ctx, &relation.RelationFollowListRequest{
		UserId: myId,
		Token:  c.Query("token"),
	})
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusInternalServerError), c, utils.H{
			"status_code": -1, //res.StatusCode
			"status_msg":  response.StatusMsg,
		})
		return
	}

	myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
		"status_code": 0, //res.StatusCode
		"status_msg":  response.StatusMsg,
		"user_list":   response.UserList,
	})
}
func RelationAction(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get(viper.Read("jwt").GetString("IdentityKey"))
	toUIdstring := c.Query("to_user_id")
	toUId, err := strconv.ParseInt(toUIdstring, 10, 64)
	myId := int64(v.(*model.User).ID)
	if myId == toUId {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "不能关注自己",
		})
		return
	}
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "to_user_id格式错误",
		})
		return
	}
	actTypeString := c.Query("action_type")
	actType, err := strconv.ParseInt(actTypeString, 10, 32)
	if err != nil || (actType != 1 && actType != 2) {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "action_type格式错误",
		})
		return
	}
	response, err := rpc.RelationAction(ctx, &relation.RelationActionRequest{
		Token:      c.Query("token"),
		ToUserId:   toUId,
		ActionType: int32(actType),
		Id:         myId,
	})
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusInternalServerError), c, utils.H{
			"status_code": -1, //res.StatusCode
			"status_msg":  response.StatusMsg,
		})
		return
	}

	myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
		"status_code": 0, //res.StatusCode
		"status_msg":  response.StatusMsg,
	})
}
