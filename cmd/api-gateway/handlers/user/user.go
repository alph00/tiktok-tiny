package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	kitex "github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func init() {
	fmt.Println("先执行init函数")
}

// Register 注册
func Register(ctx context.Context, c *app.RequestContext) {
	fmt.Println("调用了api/handler/user下的user.go register函数")
	username := c.Query("username")
	password := c.Query("password")
	//校验参数
	if len(username) == 0 || len(password) == 0 {
		myutil.ResponseMsg(http.StatusBadRequest, c, utils.H{
			"status_code": -1,
			"status_msg":  "用户名或密码不能为空", //res.StatusMsg,
		})
		return
	}
	if len(username) > 32 || len(password) > 32 {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "用户名或密码长度不能大于32个字符", //res.StatusMsg,
		})
		return
	}
	//调用kitex/kitex_gen
	fmt.Println("调用kitex/kitex_gen")
	req := &kitex.UserRegisterRequest{
		Username: username,
		Password: password,
	}
	fmt.Println("进行rpc注册，调用rpc/user下的Register函数")
	res, _ := rpc.Register(ctx, req)
	fmt.Printf("进行rpc注册，调用rpc/user下的Register函数完成,%v\n", res)
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
		"user_id":     res.UserId,
		"token":       res.Token,
	})
}

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) int64 {
	username := c.Query("username")
	password := c.Query("password")
	//校验参数
	if len(username) == 0 || len(password) == 0 {
		myutil.ResponseMsg(http.StatusBadRequest, c, utils.H{
			"status_code": -1,
			"status_msg":  "用户名或密码不能为空", //res.StatusMsg,
		})
		return -1
	}
	//调用kitex/kitex_gen
	req := &kitex.UserLoginRequest{
		Username: username,
		Password: password,
	}
	res, _ := rpc.Login(ctx, req)
	if res.StatusCode == -1 {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg, //res.StatusMsg,
		})
		return -1
	}
	myutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg, //res.StatusMsg,
		"user_id":     res.UserId,
	})
	return res.UserId
}

// UserInfo 用户信息
func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	token := c.Query("token")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  "user_id 不合法", //res.StatusMsg,
			"user":        nil,
		})
		return
	}

	//调用kitex/kitex_genit
	req := &kitex.UserInfoRequest{
		UserId: id,
		Token:  token,
	}
	res, _ := rpc.UserInfo(ctx, req)
	if res.StatusCode == -1 {
		myutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg, //res.StatusMsg,
			"user":        nil,
		})
		return
	}
	myutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg, //res.StatusMsg,
		"user":        res.User,
	})
}
