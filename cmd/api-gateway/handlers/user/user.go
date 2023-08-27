package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	rpc "github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc/user"
	kitex "github.com/alph00/tiktok-tiny/kitex_gen/user"

	"github.com/alph00/tiktok-tiny/internal/response"
	"github.com/cloudwego/hertz/pkg/app"
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
		c.JSON(http.StatusBadRequest, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	if len(username) > 32 || len(password) > 32 {
		c.JSON(http.StatusOK, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码长度不能大于32个字符",
			},
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
		c.JSON(http.StatusOK, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Register{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println("调用了api/handler/user下的user.go Login函数")
	//校验参数
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	//调用kitex/kitex_gen
	req := &kitex.UserLoginRequest{
		Username: username,
		Password: password,
	}
	fmt.Println("进行rpc注册，调用rpc/user下的Login函数")
	res, _ := rpc.Login(ctx, req)
	fmt.Printf("进行rpc注册，调用rpc/user下的Login函数完成,%v\n", res)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Login{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// UserInfo 用户信息
func UserInfo(ctx context.Context, c *app.RequestContext) {
	fmt.Println("调用了api/handler/user下的user.go UserInfo函数")
	userId := c.Query("user_id")
	token := c.Query("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "token 已过期",
			},
			User: nil,
		})
		return
	}
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
			User: nil,
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
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			User: nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.UserInfo{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		User: res.User,
	})
}
