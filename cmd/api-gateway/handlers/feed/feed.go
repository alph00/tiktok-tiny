package feed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jsonutil "github.com/alph00/tiktok-tiny/cmd/api-gateway/handlers/myutil"
	"github.com/alph00/tiktok-tiny/cmd/api-gateway/rpc"
	"github.com/alph00/tiktok-tiny/kitex_gen/feed"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	jwt "github.com/golang-jwt/jwt/v4"
	hretz_jwt "github.com/hertz-contrib/jwt"
)

func init() {
	fmt.Println("-----------init------------")
}

// TODO 视频流接口
func Feed(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	latestTime := c.Query("latest_time")
	var timestamp int64 = 0
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}
	var myId int64
	var req *feed.FeedRequest

	if token == "" {
		req = &feed.FeedRequest{
			LatestTime: timestamp,
			Token:      token,
			HasLogin:   false,
		}
		fmt.Printf("myId: %v\n", myId)
	} else {
		getIdentity, _ := hretz_jwt.New(&hretz_jwt.HertzJWTMiddleware{
			Realm:         "test zone",
			Key:           []byte(viper.Read("jwt").GetString("Key")),
			Timeout:       time.Hour,
			MaxRefresh:    time.Hour,
			TokenLookup:   "query: token, header: Authorization, cookie: jwt",
			TokenHeadName: "Bearer",
			IdentityKey:   viper.Read("jwt").GetString("IdentityKey"),
		})
		// fmt.Printf("token: %v\n", token)
		res, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod(getIdentity.SigningAlgorithm) != t.Method {
				return nil, errors.New("invalid signing algorithm")
			}
			if true {
				return getIdentity.Key, nil
			}
			// save token string if valid
			c.Set("JWT_TOKEN", token)

			return getIdentity.Key, nil
		}, getIdentity.ParseOptions...)

		myId := res.Claims.(jwt.MapClaims)[getIdentity.IdentityKey]
		// fmt.Printf("token: %v\n", token)
		// fmt.Printf("myId: %v\n", myId)
		req = &feed.FeedRequest{
			LatestTime: timestamp,
			Token:      token,
			Id:         int64(myId.(float64)),
			HasLogin:   true,
		}
	}
	res, _ := rpc.Feed(ctx, req)
	if res.StatusCode == -1 {
		jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
			"status_code": -1,
			"status_msg":  res.StatusMsg,
			"video_list":  nil,
		})
		return
	}

	jsonutil.ResponseMsg(http.StatusOK, c, utils.H{
		"status_code": 0,
		"status_msg":  res.StatusMsg,
		"video_list":  res.VideoList,
	})

}
