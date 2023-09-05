package service

import (
	"context"

	"github.com/alph00/tiktok-tiny/kitex_gen/comment"
	userInfo "github.com/alph00/tiktok-tiny/kitex_gen/user"
	"github.com/alph00/tiktok-tiny/model"
	"github.com/alph00/tiktok-tiny/pkg/minio"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
// TODO 评论操作 登录用户对视频进行评论 /douyin/comment/action
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {

	//TODO 评论操作   1-发布评论,2-删除评论
	actionType := req.ActionType
	//TODO  根据视频id获取视频信息
	video, err := model.GetVideoInfoById(ctx, req.VideoId)
	if err != nil {
		return &comment.CommentActionResponse{
			StatusMsg:  "该id对应的视频不存在",
			StatusCode: -1,
		}, nil
	}
	//TODO 判断评论行为  1-发布评论
	//TODO  获取用户id
	userId := req.Id
	if actionType == 1 {
		//TODO 实例化 评论模型
		commentModel := &model.Comment{
			VideoID: uint(req.VideoId),
			UserID:  uint(userId),
			Content: req.CommentText,
		}
		//TODO 把评论信息写入数据库中
		err := model.InsertComment(ctx, commentModel)
		if err != nil {
			res := &comment.CommentActionResponse{
				StatusMsg:  "评论消息保存失败",
				StatusCode: -1,
			}
			return res, nil
		}
		//TODO 	2-删除评论
	} else if actionType == 2 {
		//TODO 评论可以删除的只有是登录用户发表的且存在的评论或者视频发表者发表的视频下的评论
		//TODO 1.根据评论id 查询该评论是否存在
		resComment, err := model.SelectCommentByCommentID(ctx, req.CommentId)
		if err != nil {
			res := &comment.CommentActionResponse{
				StatusMsg:  "获取评论消息失败",
				StatusCode: -1,
			}
			return res, nil
		}
		//TODO 2.判断查询的评论所属的用户信息,用户信息与登录用户的信息匹配,才可以删除对应的评论
		if resComment == nil {
			res := &comment.CommentActionResponse{
				StatusMsg:  "该评论不存在，无法删除",
				StatusCode: -1,
			}
			return res, nil
		} else {
			//TODO 获取该评论对应的视频信息
			videoByComment, err := model.GetVideoInfoById(ctx, int64(resComment.VideoID))
			if err != nil {
				res := &comment.CommentActionResponse{
					StatusMsg:  "系统出现异常，评论无法删除",
					StatusCode: -1,
				}
				return res, nil
			}
			//TODO 如果用户不是评论的发表评论的用户或者视频的发表者无法删除评论
			if userId != int64(resComment.UserID) || userId != int64(videoByComment.AuthorID) {
				res := &comment.CommentActionResponse{
					StatusMsg:  "用户无删除权限，评论无法删除",
					StatusCode: -1,
				}
				return res, nil
			}
		}
		//TODO 删除评论操作
		err = model.DelCommentByID(ctx, int64(resComment.ID), int64(video.ID))
		if err != nil {
			res := &comment.CommentActionResponse{
				StatusMsg:  "服务内部异常，评论无法删除",
				StatusCode: -1,
			}
			return res, nil
		}

	} else {
		//TODO
		res := &comment.CommentActionResponse{
			StatusMsg:  "action_type 字段信息有误",
			StatusCode: -1,
		}
		return res, nil
	}

	//TODO 评论删除成功
	res := &comment.CommentActionResponse{
		StatusMsg:  "success",
		StatusCode: 0,
		Comment: &comment.Comment{
			Content: req.CommentText,
		},
	}
	return res, nil

}

// CommentList implements the CommentServiceImpl interface.
// TODO 查看视频的所有的评论 按照发布时间进行排序 /douyin/comment/list 评论列表
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {

	//TODO   根据视频id获取指定视频的全部评论内容
	commentList, err := model.SelectVideoCommentListByVideoID(ctx, req.VideoId)
	if err != nil {
		res := &comment.CommentListResponse{
			StatusMsg:  "服务异常,获取评论信息失败",
			StatusCode: -1,
		}
		return res, nil
	}
	//TODO 对于评论的查询结果进行封装与字段映射得到查询的评论列表的结果
	comments := make([]*comment.Comment, 0)
	//TODO for 循环遍历
	for _, c := range commentList {
		//TODO 根据评论对应的用户id 查询用户信息
		user, err := model.GetUserByID(ctx, int64(c.UserID))
		if err != nil {
			res := &comment.CommentListResponse{
				StatusMsg:  "服务异常,获取评论信息失败",
				StatusCode: -1,
			}
			return res, nil
		}
		//TODO  获取用户id

		userId := req.Id
		//TODO  根据id获取用户之间的关注关系
		isFlolow := model.IsFollow(int64(user.ID), int64(userId))

		//TODO 从minio中获取头像URL
		avatar, err := minio.GetFileTemporaryURL(minio.Avatar, user.Avatar)
		if err != nil {
			res := &comment.CommentListResponse{
				StatusMsg:  "服务异常,获取头像信息失败",
				StatusCode: -1,
			}
			return res, nil
		}
		//TODO 从minio中获取背景图片URL
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImage, user.BackgroundImage)
		if err != nil {
			res := &comment.CommentListResponse{
				StatusMsg:  "服务异常,获取背景图片信息失败",
				StatusCode: -1,
			}
			return res, nil
		}
		//TODO user 模型对象封装
		usrInfo := &userInfo.User{
			Id:              int64(userId),
			Name:            user.UserName,
			FollowCount:     int64(user.FollowingCount),
			FollowerCount:   int64(user.FollowerCount),
			IsFollow:        isFlolow,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       user.Signature,
			TotalFavorited:  int64(user.TotalFavorited),
			WorkCount:       int64(user.WorkCount),
			FavoriteCount:   int64(user.FavoriteCount),
		}
		//TODO Comment 模型对象封装 添加到comments的切片中
		if int64(c.ID) > 0 {
			comments = append(comments, &comment.Comment{
				Id:         int64(c.ID),
				User:       usrInfo,
				Content:    c.Content,
				CreateDate: c.CreatedAt.Format("2023-08-19"),
				LikeCount:  int64(c.LikeCount),
				TeaseCount: int64(c.TeaseCount),
			})
		}

	}
	//TODO 将封装的评论列表结果返回
	res := &comment.CommentListResponse{
		StatusMsg:   "success",
		StatusCode:  0,
		CommentList: comments,
	}
	return res, nil

}
