package service

import (
	"context"
	"fmt"
	"strings"

	favorite "github.com/alph00/tiktok-tiny/kitex_gen/favorite"
	"github.com/alph00/tiktok-tiny/model"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	UserId := req.UserId
	VideoID := req.VideoId
	actType := req.ActionType
	fmt.Println("调用favorite微服务handler下的FavoriteAction函数")
	if actType == 1 {
		fmt.Println("acttype=1,来创建点赞数据")
		err := model.CreateVideoFavorite(ctx, &model.FavoriteVideoRelation{
			UserID:  uint(UserId),
			VideoID: uint(VideoID),
		})
		if err != nil {
			res := &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "向数据库创建点赞视频数据failure",
			}
			return res, nil
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
