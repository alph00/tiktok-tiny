package message

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	"github.com/alph00/tiktok-tiny/kitex_gen/message"
	"github.com/alph00/tiktok-tiny/model"
	mw "github.com/alph00/tiktok-tiny/pkg/mw/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func MessageAction(ctx context.Context, c *app.RequestContext) {
	toUIdstring := c.Query("to_user_id")
	if toUIdstring == "" {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "to_user_id为空", //res.StatusMsg,
		})
	}
	toUId, err := strconv.ParseInt(toUIdstring, 10, 64)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "to_user_id 不是int类型", //res.StatusMsg,
		})
	}
	actType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || actType != 1 {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "type 非法", //res.StatusMsg,
		})
	}
	if len(c.Query("content")) == 0 {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  "content为空", //res.StatusMsg,
		})
	}

	v, _ := c.Get(mw.IdentityKey)
	response, _ := rpc.MessageAction(ctx, &message.MessageActionRequest{
		Id:         int64(v.(*model.User).ID),
		Token:      c.Query("token"),
		ToUserId:   toUId,
		ActionType: int32(actType),
		Content:    c.Query("content"),
	})
	// if err != nil {
	// 	// to do
	// }
	if response.StatusCode == 0 {
		myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
			"status_code": 0,
			"status_msg":  response.StatusMsg, //res.StatusMsg,
		})
	} else {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code": -1,
			"status_msg":  response.StatusMsg, //res.StatusMsg,
		})
	}
}

func MessageChat(ctx context.Context, c *app.RequestContext) {
	toUIdstring := c.Query("to_user_id")
	if toUIdstring == "" {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code":  -1,
			"status_msg":   "to_user_id为空", //res.StatusMsg,
			"message_list": nil,
		})
	}
	fmt.Printf("toUIdstring: %v\n", toUIdstring)
	toUId, err := strconv.ParseInt(toUIdstring, 10, 64)
	if err != nil {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code":  -1,
			"status_msg":   "to_user_id 不是int类型", //res.StatusMsg,
			"message_list": nil,
		})
	}

	// fmt.Printf("mw.IdentityKey: %v\n", mw.IdentityKey)
	v, _ := c.Get(mw.IdentityKey)
	// debug := int64(v.(*model.User).ID)
	// fmt.Printf("debug: %v\n", debug)
	response, err := rpc.MessageChat(ctx, &message.MessageChatRequest{
		Id:       int64(v.(*model.User).ID),
		Token:    c.Query("token"),
		ToUserId: toUId,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if response.StatusCode == 0 {
		myutil.ResponseMsg(int(consts.StatusOK), c, utils.H{
			"status_code":  0,
			"status_msg":   response.StatusMsg, //res.StatusMsg,
			"message_list": response.MessageList,
		})
	} else {
		myutil.ResponseMsg(int(consts.StatusBadRequest), c, utils.H{
			"status_code":  -1,
			"status_msg":   response.StatusMsg, //res.StatusMsg,
			"message_list": nil,
		})
	}
}
