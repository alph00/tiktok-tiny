package main

import (
	"context"
	"time"

	auth "github.com/alph00/tiktok-tiny/kitex_gen/auth"
	"github.com/alph00/tiktok-tiny/model/auth/db"
	"github.com/alph00/tiktok-tiny/package/jwt"
	"github.com/bytedance/gopkg/util/logger"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}

// Register implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Register(ctx context.Context, req *auth.UserRegisterRequest) (resp *auth.UserRegisterResponse, err error) {
	// TODO: Your code here...
	// TODO:日志打印

	// 检查用户名是否冲突
	usr, err := db.GetUserByName(ctx, req.Username)
	if err != nil {
		// logger.Errorln(err.Error())
		res := &auth.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "注册失败：服务器内部错误",
		}
		return res, nil
	} else if usr != nil {
		// logger.Errorf("该用户名已存在：%s", usr.UserName)
		res := &auth.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "该用户名已存在，请更换",
		}
		return res, nil
	}

	// 创建user
	// rand.Seed(time.Now().UnixMilli())
	usr = &db.Auth{
		UserName: req.Username,
		// TODO:密码加密
		Password: req.Password,
		// Avatar:   fmt.Sprintf("default%d.png", rand.Intn(10)),
	}
	if err := db.CreateUser(ctx, usr); err != nil {
		// TODO：日志采集
		// logger.Errorln(err.Error())
		res := &auth.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "注册失败：服务器内部错误",
		}
		return res, nil
	}
	//生成token
	claims := jwt.CustomClaims{Id: int64(usr.ID)}
	claims.ExpiresAt = time.Now().Add(time.Minute * 5).Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
		res := &auth.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：token 创建失败",
		}
		return res, nil
	}
	res := &auth.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}

	return res, nil
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.UserLoginRequest) (resp *auth.UserLoginResponse, err error) {
	// TODO: Your code here...
	// TODO:日志打印

	// 根据用户名获取密码
	usr, err := db.GetUserByName(ctx, req.Username)
	if err != nil {
		// logger.Errorln(err.Error())
		res := &auth.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "登录失败：服务器内部错误",
		}
		return res, nil
	} else if usr == nil {
		res := &auth.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名不存在",
		}
		return res, nil
	}

	// 比较数据库中的密码和请求的密码
	// TODO:密码加密处理
	if req.Password != usr.Password {
		logger.Error("用户名或密码错误")
		res := &auth.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名或密码错误",
		}
		return res, nil
	}

	// 密码认证通过,获取用户id并生成token
	claims := jwt.CustomClaims{
		Id: int64(usr.ID),
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
		res := &auth.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：token 创建失败",
		}
		return res, nil
	}

	// 返回结果
	res := &auth.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return res, nil

	// return
}
