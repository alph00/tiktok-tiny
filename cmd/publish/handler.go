package main

import (
	"context"
	"fmt"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/minio"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	"time"
)

// PublishServiceImpl implements the last service interface defined in the IDL.
type PublishServiceImpl struct{}

// PublishAction implements the PublishServiceImpl interface.
// TODO 登录用户选择视频上传业务逻辑  /douyin/publish/action
func (s *PublishServiceImpl) PublishAction(ctx context.Context, req *publish.PublishActionRequest) (resp *publish.PublishActionResponse, err error) {
	// TODO: 视频投稿功能实现
	title := req.Title
	token := req.Token
	data := req.Data
	var userId int64 = -1
	//TODO 校验token的有效性
	fmt.Println(token)
	if token != "" {
		claims, err := Jwt.ParseToken(token)
		fmt.Println(err)
		fmt.Println(claims)
		//TODO 如果存在异常，token校验失败
		if err != nil {
			res := &publish.PublishActionResponse{
				StatusCode: -1,
				StatusMsg:  "token解析出现异常",
			}
			return res, nil
		}
		//TODO 获取到登录用户id
		userId = claims.Id
	}
	//TODO 校验视频标题的长度
	if len(title) == 0 || len(title) > 32 {
		res := &publish.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "标题不能为空同时不能超过32个字符",
		}
		return res, nil
	}
	//TODO 校验文件的大小
	maxSize := viper.Init("video").Viper.GetInt("video.maxSizeLimit")
	if len(data) > maxSize*1000*1000 {
		res := &publish.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("视频大小为%v,超出%vM要求", len(data), maxSize),
		}
		return res, nil
	}
	//TODO 数据库信息收集
	createTimestamp := time.Now().UnixMilli()
	videoTitle := fmt.Sprintf("%d_%s_%d.mp4", userId, title, createTimestamp)
	coverTitle := fmt.Sprintf("%d_%s_%d.png", userId, title, createTimestamp)
	//TODO 数据库插入操作实体类封装
	v := &model.Video{
		Title:       title,
		PlayUrl:     videoTitle,
		CoverUrl:    coverTitle,
		AuthorID:    uint(userId),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	//TODO 视频信息保存
	err = model.CreateVideo(ctx, v)
	fmt.Println(err)
	if err != nil {
		res := &publish.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频发布失败，服务内部异常",
		}
		return res, nil
	}
	//TODO   利用协程异步 上传视频至minio
	go func() {
		//TODO 视频上传到 minio
		err := VideoPublish(req.Data, videoTitle, coverTitle)
		if err != nil {
			// TODO 出现异常 删除插入的记录
			e := model.DelVideoByID(ctx, int64(v.ID), userId)
			if e != nil {
				logger.Errorf("视频记录删除失败：%s", err.Error())
			}
		}
	}()
	//TODO 视频上传成功
	res := &publish.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "视频创建成功，等待后台上传完成",
	}
	return res, nil

}

// PublishList implements the PublishServiceImpl interface.
// TODO 登录用户的视频发布列表,直接列出用于所有投稿过的视频 /douyin/publish/list
func (s *PublishServiceImpl) PublishList(ctx context.Context, req *publish.PublishListRequest) (resp *publish.PublishListResponse, err error) {
	// TODO: 获取登录用户的id
	userId := req.UserId
	//TODO 根据用户id 查询数据库获取用户的所有投稿视频信息
	resultList, err := model.GetVideoListByUserID(ctx, userId)
	if err != nil {
		logger.Infof("查询数据库出现异常,获取视频列表信息失败")
		res := &publish.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  "查询数据库出现异常,获取视频列表信息失败",
		}
		return res, nil
	}
	//TODO 对于数据库的查询结果进行遍历,并将查询结果进行转化以及字段映射,并放入切片中
	videoList := make([]*publish.Video, 0)
	//TODO  for 循环遍历
	for _, video := range resultList {
		//TODO 查询出关注的用户信息
		authInfo, err := model.GetUserByID(ctx, int64(video.AuthorID))
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取用户信息失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 根据用户id 获取出当前用户的关注用户 "user_id=? AND to_user_id=?
		follow, err := model.GetRelationByUserIDs(ctx, userId, int64(authInfo.ID))
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取关注信息失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 从用户的点赞关系表中获取对应的点赞关系        "user_id = ? and video_id = ?
		favorite, err := model.GetFavoriteVideoRelationByUserVideoID(ctx, userId, int64(video.ID))
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取点赞关系失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 从minio中获取上传视频的URL
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, video.PlayUrl)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取视频URL失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 从minio中获取上传封面的URL
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, video.CoverUrl)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取封面URL失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 从minio中获取上传头像的URL
		avatarUrl, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, authInfo.Avatar)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取头像URL失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 从minio中获取背景图片的URL
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, authInfo.BackgroundImage)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取背景图片URL失败,服务内部出现异常",
			}
			return res, nil
		}
		//TODO 属性填充 字段映射 每个视频信息封装好后放在切片中
		videoList = append(videoList, &publish.Video{
			Id: int64(video.ID),
			Author: &publish.User{
				Id:              int64(authInfo.ID),
				Name:            authInfo.UserName,
				FollowerCount:   int64(authInfo.FollowerCount),
				FollowCount:     int64(authInfo.FollowingCount),
				IsFollow:        follow != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       authInfo.Signature,
				TotalFavorited:  int64(authInfo.TotalFavorited),
				WorkCount:       int64(authInfo.WorkCount),
				FavoriteCount:   int64(authInfo.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         video.Title,
		})

	}

	//TODO 查询以及封装结果返回
	res := &publish.PublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}

	return res, nil

}
