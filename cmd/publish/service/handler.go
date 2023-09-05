package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alph00/tiktok-tiny/kitex_gen/feed"
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
	"github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/minio"
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/bytedance/gopkg/util/logger"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PublishServiceImpl struct{}

// TODO 登录用户选择视频上传业务逻辑  /douyin/publish/action
func (s *PublishServiceImpl) PublishAction(ctx context.Context, req *publish.PublishActionRequest) (resp *publish.PublishActionResponse, err error) {
	// TODO: 视频投稿功能实现
	title, data := req.Title, req.Data
	//TODO 获取用户id
	userId := req.Id
	//TODO 校验视频标题的长度
	if len(title) == 0 || len(title) > 32 {
		return &publish.PublishActionResponse{StatusCode: -1, StatusMsg: "the title of video should in 0~32"}, nil
	}
	//TODO 校验文件的大小
	maxSize := viper.Read("service").GetInt("publish.videoMaxSize")
	if len(data) > maxSize*1000*1000 {
		return &publish.PublishActionResponse{StatusCode: -1, StatusMsg: fmt.Sprintf("the size of video is %v, greater than %vM", len(data), maxSize)}, nil
	}
	//TODO 数据库信息收集
	createTimestamp := time.Now().UnixMilli()
	videoTitle := fmt.Sprintf("%d_%s_%d.mp4", userId, title, createTimestamp)
	coverTitle := fmt.Sprintf("%d_%s_%d.png", userId, title, createTimestamp)
	//TODO 数据库插入操作实体类封装
	v := &model.Video{
		Title:     title,
		PlayUrl:   videoTitle,
		CoverUrl:  coverTitle,
		AuthorID:  uint(userId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	//TODO 视频信息保存
	err = model.SaveVideo(ctx, v)
	if err != nil {
		return &publish.PublishActionResponse{StatusCode: -1, StatusMsg: "internal error"}, nil
	}
	//
	playUrl, err := uploadVideo(req.Data, videoTitle)
	if err != nil {
		model.DelVideoByID(ctx, int64(v.ID), int64(userId))
		return &publish.PublishActionResponse{StatusCode: -1, StatusMsg: "uploading video error"}, err
	}
	err = uploadCover(playUrl, coverTitle)
	if err != nil {
		model.DelVideoByID(ctx, int64(v.ID), int64(userId))
		return &publish.PublishActionResponse{StatusCode: -1, StatusMsg: "uploading cover error"}, err
	}
	//

	// go func() {
	// 	playUrl, err := uploadVideo(req.Data, videoTitle)
	// 	if err != nil {
	// 		model.DelVideoByID(ctx, int64(v.ID), int64(userId))
	// 		return
	// 	}
	// 	err = uploadCover(playUrl, coverTitle)
	// 	if err != nil {
	// 		model.DelVideoByID(ctx, int64(v.ID), int64(userId))
	// 		return
	// 	}
	// }()
	return &publish.PublishActionResponse{StatusCode: 0, StatusMsg: "process success, uploading video"}, nil
}

// TODO 登录用户的视频发布列表,直接列出用于所有投稿过的视频 /douyin/publish/list
func (s *PublishServiceImpl) PublishList(ctx context.Context, req *publish.PublishListRequest) (resp *publish.PublishListResponse, err error) {
	// TODO: 获取登录用户的id
	userId := req.UserId
	//TODO 根据用户id 查询数据库获取用户的所有投稿视频信息
	resultList, err := model.GetVideoListByUserID(ctx, userId)
	if err != nil {
		res := &publish.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  "internal error",
		}
		return res, nil
	}
	//TODO 对于数据库的查询结果进行遍历,并将查询结果进行转化以及字段映射,并放入切片中
	videoList := make([]*feed.Video, 0)
	//TODO  for 循环遍历
	for _, video := range resultList {
		//TODO 查询出关注的用户信息
		authInfo, err := model.GetUserByID(ctx, int64(video.AuthorID))
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取用户信息失败,500 error",
			}
			return res, nil
		}
		isFollow := model.IsFollow(int64(authInfo.ID), userId)
		isFavorite := model.IsFavorite(userId, int64(video.ID))
		//TODO 从minio中获取上传视频的URL
		playUrl, err := minio.GetFileTemporaryURL(minio.Video, video.PlayUrl)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取视频URL失败,500 error",
			}
			return res, nil
		}
		//TODO 从minio中获取上传封面的URL
		coverUrl, err := minio.GetFileTemporaryURL(minio.Cover, video.CoverUrl)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取封面URL失败,500 error",
			}
			return res, nil
		}
		//TODO 从minio中获取上传头像的URL
		avatarUrl, err := minio.GetFileTemporaryURL(minio.Avatar, authInfo.Avatar)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取头像URL失败,500 error",
			}
			return res, nil
		}
		//TODO 从minio中获取背景图片的URL
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImage, authInfo.BackgroundImage)
		if err != nil {
			res := &publish.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "获取背景图片URL失败,500 error",
			}
			return res, nil
		}
		//TODO 属性填充 字段映射 每个视频信息封装好后放在切片中
		videoList = append(videoList, &feed.Video{
			Id: int64(video.ID),
			Author: &user.User{
				Id:              int64(authInfo.ID),
				Name:            authInfo.UserName,
				FollowerCount:   int64(authInfo.FollowerCount),
				FollowCount:     int64(authInfo.FollowingCount),
				IsFollow:        isFollow,
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
			IsFavorite:    isFavorite,
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

func uploadVideo(data []byte, title string) (string, error) {

	//TODO 读取视频的文件流
	filereader := bytes.NewReader(data)
	//TODO 文件的类型
	contentType := "application/mp4"
	//TODO minio 文件上传
	uploadSize, err := minio.UploadFileByIO(minio.Video, title, filereader, int64(len(data)), contentType)
	if err != nil {
		return "", err
	}
	logger.Infof("视频文件大小为：%v", uploadSize)
	//TODO 获取minio中的文件路径
	playUrl, err := minio.GetFileTemporaryURL(minio.Video, title)
	if err != nil {
		return "", err
	}
	return playUrl, nil
}

// TODO 上传封面至 minio
func uploadCover(playUrl string, title string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(playUrl).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err
	}
	//TODO 封面图片字节流写入
	var imgByte []byte
	buf.Write(imgByte)
	contentType := "image/png"
	//TODO 封面上传至minio
	_, err = minio.UploadFileByIO(minio.Cover, title, buf, int64(buf.Len()), contentType)
	if err != nil {
		return err
	}
	return nil
}
