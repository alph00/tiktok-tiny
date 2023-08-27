package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	user "github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/tools"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang-jwt/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	// 检查用户名是否冲突
	usr, err := model.GetUserByName(ctx, req.Username)
	if err != nil {
		res := &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "注册失败：服务器内部错误",
		}
		return res, nil
	} else if usr != nil {
		res := &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "该用户名已存在，请更换",
		}
		return res, nil
	}

	// 创建user
	rand.Seed(time.Now().UnixMilli())
	usr = &model.User{
		UserName: req.Username,
		Password: tools.Md5Encrypt(req.Password),
		Avatar:   fmt.Sprintf("default%d.png", rand.Intn(10)),
	}
	if err := model.CreateUser(ctx, usr); err != nil {
		res := &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "注册失败：服务器内部错误",
		}
		return res, nil
	}

	//if err := db.CreateUsers(ctx, []*db.User{{
	//	UserName: req.Username,
	//	Password: tool.Md5Encrypt(req.Password),
	//	Avatar:   fmt.Sprintf("default%d.png", rand.Intn(10)),
	//}}); err != nil {
	//	logger.Errorf("发生错误：%v", err.Error())
	//	res := &user.UserRegisterResponse{
	//		StatusCode: -1,
	//		StatusMsg:  "注册失败：服务器内部错误",
	//	}
	//	return res, nil
	//}

	// 获取用户id
	//usr, err = db.GetUserByName(ctx, req.Username)
	//if err != nil || usr == nil {
	//	logger.Errorf("发生错误：%v", err.Error())
	//	res := &user.UserRegisterResponse{
	//		StatusCode: -1,
	//		StatusMsg:  "注册失败：服务器内部错误",
	//	}
	//	return res, nil
	//}

	// Create the token
	//这里目前只能支持公玥签发
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["Id"] = int64(usr.ID)
	expire := time.Now().Add(time.Hour)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()
	tokenString, err := token.SignedString([]byte("secret key"))

	if err != nil {
		res := &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：token 创建失败",
		}
		return res, nil
	}
	res := &user.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      tokenString,
	}
	return res, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	// 根据用户名获取密码
	usr, err := model.GetUserByName(ctx, req.Username)
	if err != nil {
		res := &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "登录失败：服务器内部错误",
		}
		return res, nil
	} else if usr == nil {
		res := &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名不存在",
		}
		return res, nil
	}

	// 比较数据库中的密码和请求的密码
	if tools.Md5Encrypt(req.Password) != usr.Password {
		res := &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名或密码错误",
		}
		return res, nil
	}

	// 返回结果
	res := &user.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
	}
	return res, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	userID := req.UserId

	// 从数据库获取user
	usr, err := model.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
		res := &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}
		return res, nil
	} else if usr == nil {
		logger.Errorf("该用户不存在：%v", err.Error())
		res := &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "该用户不存在",
		}
		return res, nil
	}

	// avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
	// if err != nil {
	// 	logger.Errorf("Minio获取头像失败：%v", err.Error())
	// 	res := &user.UserInfoResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "服务器内部错误：获取头像失败",
	// 	}
	// 	return res, nil
	// }
	// backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, usr.BackgroundImage)
	// if err != nil {
	// 	logger.Errorf("Minio获取背景图失败：%v", err.Error())
	// 	res := &user.UserInfoResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "服务器内部错误：获取背景图失败",
	// 	}
	// 	return res, nil
	// }

	//返回结果
	res := &user.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User: &user.User{
			Id:            int64(usr.ID),
			Name:          usr.UserName,
			FollowCount:   int64(usr.FollowingCount),
			FollowerCount: int64(usr.FollowerCount),
			IsFollow:      userID == int64(usr.ID),
			// Avatar:          avatar,
			// BackgroundImage: backgroundImage,
			Signature:      usr.Signature,
			TotalFavorited: int64(usr.TotalFavorited),
			WorkCount:      int64(usr.WorkCount),
			FavoriteCount:  int64(usr.FavoriteCount),
		},
	}
	return res, nil
}
