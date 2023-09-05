package service

import (
	"context"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/feed"
	"github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/minio"
)

type FeedServiceImpl struct{}

func (s *FeedServiceImpl) Feed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	nextTime := time.Now().UnixMilli()

	//TODO 查找最近30条的视频
	videos, err := model.GetLatestTimeVideos(ctx, 2, &req.LatestTime)
	if err != nil {
		return &feed.FeedResponse{StatusCode: -1, StatusMsg: "internal error"}, nil
	}

	//TODO 视频列表的字段填充
	videoList := make([]*feed.Video, 0)
	for _, video := range videos {
		author, err := model.GetUserByID(ctx, int64(video.AuthorID))
		if err != nil {
			return nil, err
		}
		var isFollow, isFavorite bool
		if req.HasLogin {
			userId := req.Id
			isFollow = model.IsFollow(int64(video.AuthorID), int64(userId))
			isFavorite = model.IsFavorite(int64(userId), int64(video.ID))
		} else {
			isFollow = false
			isFavorite = false
		}

		playUrl, err := minio.GetFileTemporaryURL(minio.Video, video.PlayUrl)
		if err != nil {
			return &feed.FeedResponse{StatusCode: -1, StatusMsg: "internal error：视频获取失败"}, nil
		}
		coverUrl, err := minio.GetFileTemporaryURL(minio.Cover, video.CoverUrl)
		if err != nil {
			return &feed.FeedResponse{StatusCode: -1, StatusMsg: "internal error：封面获取失败"}, nil
		}
		avatarUrl, err := minio.GetFileTemporaryURL(minio.Avatar, author.Avatar)
		if err != nil {
			return &feed.FeedResponse{StatusCode: -1, StatusMsg: "internal error：头像获取失败"}, nil
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImage, author.BackgroundImage)
		if err != nil {
			return &feed.FeedResponse{StatusCode: -1, StatusMsg: "internal error：背景图获取失败"}, nil
		}
		videoList = append(videoList, &feed.Video{
			Id: int64(video.ID),
			Author: &user.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowingCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        isFollow,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	return &feed.FeedResponse{StatusCode: 0, StatusMsg: "success", VideoList: videoList, NextTime: nextTime}, nil
}
