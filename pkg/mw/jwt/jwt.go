/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/user"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	JwtConfig     = viper.Read("jwt")
	IdentityKey   = JwtConfig.GetString("IdentityKey")
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		//能返回jwtcookie吗？
		SendCookie:    true,
		Realm:         "test zone",
		Key:           []byte(JwtConfig.GetString("Key")),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "query: token, header: Authorization, cookie: jwt",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			// c.JSON(http.StatusOK, utils.H{
			// 	// "code":    code,
			// 	"token": token,
			// 	// "expire":  expire.Format(time.RFC3339),
			// 	// "message": "success",
			// })
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				// fmt.Printf("v: %v\n", v)
				return jwt.MapClaims{
					IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			res := user.Login(ctx, c)
			if res == -1 {
				return nil, errors.New("登录失败")
			} else {
				return res, nil
			}
			//以下为暂时替代login
			// if true {
			// 	c.JSON(http.StatusOK, utils.H{
			// 		"status_code": 0,
			// 		"status_msg":  "登录成功", //res.StatusMsg,
			// 		"user_id":     1,      //res.UserId,
			// 	})
			// } else {
			// 	c.JSON(http.StatusOK, utils.H{
			// 		"status_code": -1,
			// 		"status_msg":  "登录失败", //res.StatusMsg,
			// 	})
			// }

			// //以下为暂时替代login
			// return int64(1), nil
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			// fmt.Printf("claims: %v\n", claims)
			res := &model.User{}
			res.ID = uint(claims[IdentityKey].(float64))
			return res
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			fmt.Printf("message: %v\n", message)
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
