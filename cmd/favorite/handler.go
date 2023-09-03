package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	favorite "github.com/alph00/tiktok-tiny/kitex_gen/favorite"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/rabbitmq"
	"github.com/alph00/tiktok-tiny/pkg/viper"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

var (
	config     = viper.Read("rabbitmq")
	autoAck    = config.GetBool("consumer.favorite.autoAck")
	FavoriteMq = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
	// err        error
)

func init() {
	go consume()
}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	// TODO: Your code here...
	// 这里的设置是req为
	// UserId:     int64(v.(*model.User).ID),
	//  Token:      token,
	// 	VideoId:    id,
	// 	ActionType: int32(acttype),
	// 通过token得到点赞人的数据ID，只对Video做点赞处理，那么对于作者的  ‘获赞总数’ 还需要保证数据库一致性，之后数据访问还需要来查找Video
	// 需要同时修改两个数据库  users和videos   A----->B.video    A.favotite_count++    B.video.favotite_count++
	UserId := req.UserId
	VideoID := req.VideoId
	actType := req.ActionType
	fmt.Println("调用favorite微服务handler下的FavoriteAction函数")
	if actType == 1 {
		fmt.Println("acttype=1,来创建点赞数据")
		// 这里不直接修改数据库，而是传递给消息队列，等消费
		// err := model.CreateVideoFavorite(ctx, &model.FavoriteVideoRelation{
		// 	UserID:  uint(UserId),
		// 	VideoID: uint(VideoID),
		// })
		message := model.FavoriteVideoRelation{
			UserID:  uint(UserId),
			VideoID: uint(VideoID),
		}

		jsonFC, _ := json.Marshal(message)
		fmt.Println("Publish new message: ", message)
		// FavoriteMq := rabbitmq.NewRabbitMQSimple("favorite", true)
		if err = FavoriteMq.PublishSimple(ctx, jsonFC); err != nil {
			log.Printf("消息队列发布错误：%v", err.Error())
			if strings.Contains(err.Error(), "连接断开") {
				// 检测到通道关闭，则重连
				go FavoriteMq.Destroy()
				FavoriteMq = rabbitmq.NewRabbitMQSimple("favorite", true)
				// 这里应该再重传一次吧TODO
				log.Print("消息队列通道关闭，正在重连")
				go consume()
				res := &favorite.FavoriteActionResponse{
					StatusCode: 0,
					StatusMsg:  "success",
				}
				return res, nil
			} else {
				res := &favorite.FavoriteActionResponse{
					StatusCode: -1,
					StatusMsg:  "向数据库创建点赞视频数据failure",
				}
				return res, nil
			}
		}
	} else {
		fmt.Println("acttype=0,来取消点赞数据")
		err := model.DelVideoFavorite(ctx, &model.FavoriteVideoRelation{
			UserID:  uint(UserId),
			VideoID: uint(VideoID),
		})
		if err != nil {
			res := &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "向数据库删除点赞视频数据failure",
			}
			return res, nil
		}
		// res := &favorite.FavoriteActionResponse{
		// 	StatusCode: -1,
		// 	StatusMsg:  "取消点赞failure",
		// }
		// return res, nil
	}
	res := &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return res, nil

	// return
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	// TODO: Your code here...
	// return
	// 这里的设置是req为
	// UserId:     int64(v.(*model.User).ID),
	//  Token:      token,
	// UserId := req.UserId
	// // VideoID := req.VideoId
	fmt.Println("调用favorite微服务下的FavoriteList函数")
	videos, err := model.ShowVideoFavorite(ctx, uint(req.UserId))
	if err != nil {
		res := &favorite.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  "数据库查询点赞视频数据failure",
			VideoList:  nil,
		}
		return res, nil
	}
	// if actType == 1 {
	// 	err := model.CreateVideoFavorite(ctx, &model.FavoriteVideoRelation{
	// 		UserID:  uint(UserId),
	// 		VideoID: uint(VideoID),
	// 	})
	// 	if err != nil {
	// 		res := &favorite.FavoriteActionResponse{
	// 			StatusCode: -1,
	// 			StatusMsg:  "向数据库创建点赞视频数据failure",
	// 		}
	// 		return res, nil
	// 	}
	// } else {
	// 	err := model.DelVideoFavorite(ctx, &model.FavoriteVideoRelation{
	// 		UserID:  uint(UserId),
	// 		VideoID: uint(VideoID),
	// 	})
	// 	if err != nil {
	// 		res := &favorite.FavoriteActionResponse{
	// 			StatusCode: -1,
	// 			StatusMsg:  "向数据库删除点赞视频数据failure",
	// 		}
	// 		return res, nil
	// 	}
	// 	// res := &favorite.FavoriteActionResponse{
	// 	// 	StatusCode: -1,
	// 	// 	StatusMsg:  "取消点赞failure",
	// 	// }
	// 	// return res, nil
	// }
	// 使用一个切片来存储转换后的字符串

	videoStrings := make([]string, len(videos))

	// 遍历切片中的每个 int 指针，并将其值转换为字符串
	for i, video := range videos {
		videoStrings[i] = fmt.Sprintf("%d", video)
	}

	// 使用逗号连接字符串切片中的元素
	output := "[" + strings.Join(videoStrings, ", ") + "]"
	res := &favorite.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success" + output,
		VideoList:  nil,
	}
	return res, nil
	// return
}
