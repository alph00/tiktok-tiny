package service

import (
	"context"
	"fmt"

	"github.com/alph00/tiktok-tiny/kitex_gen/user"

	relation "github.com/alph00/tiktok-tiny/kitex_gen/relation"
	"github.com/alph00/tiktok-tiny/model"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	toUId, err := model.GetUserByID(ctx, req.ToUserId)
	if err != nil {
		return &relation.RelationActionResponse{StatusCode: -1, StatusMsg: "无法获取用户信息"}, err
	}
	if toUId == nil {
		return &relation.RelationActionResponse{StatusCode: -1, StatusMsg: "该用户不存在"}, nil
	}

	if req.ActionType == 1 {
		if model.IsFollow(req.ToUserId, req.Id) {
			return &relation.RelationActionResponse{StatusCode: 0, StatusMsg: "已经关注过该用户"}, nil //这样算成功吗？
		}
		err := model.Follow(req.ToUserId, req.Id)
		if err != nil {
			return &relation.RelationActionResponse{StatusCode: -1, StatusMsg: "关注失败"}, err
		}
		model.AddFollowerCount(req.ToUserId, 1)
		model.AddFollowingCount(req.Id, 1)
	} else if req.ActionType == 2 {
		if !model.IsFollow(req.ToUserId, req.Id) {
			return &relation.RelationActionResponse{StatusCode: 0, StatusMsg: "没有关注过该用户"}, nil //这样算成功吗？
		}
		err := model.UnFollow(req.ToUserId, req.Id)
		if err != nil {
			return &relation.RelationActionResponse{StatusCode: -1, StatusMsg: "取消关注失败"}, err
		}
		model.ReduceFollowerCount(req.ToUserId, 1)
		model.ReduceFollowingCount(req.Id, 1)
	}
	return &relation.RelationActionResponse{StatusCode: 0, StatusMsg: "success"}, nil
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.RelationFollowListRequest) (resp *relation.RelationFollowListResponse, err error) {
	followeds, err := model.GetFollowedById(req.UserId)
	if err != nil {
		return &relation.RelationFollowListResponse{StatusCode: -1, StatusMsg: "获取关注列表ID失败"}, err
	}

	followedsUser, err := model.GetUsersByIDs(ctx, followeds)
	if err != nil {
		return &relation.RelationFollowListResponse{StatusCode: -1, StatusMsg: "获取关注列表USER失败"}, err
	}
	var userList []*user.User
	for _, followed := range followedsUser {
		userList = append(userList, &user.User{
			Id:            int64(followed.ID),
			Name:          followed.UserName,
			FollowCount:   int64(followed.FollowingCount),
			FollowerCount: int64(followed.FollowerCount),
			IsFollow:      true,
			// Avatar:          avatar,//to do
			// BackgroundImage: backgroundUrl,//to do
			Signature:      followed.Signature,
			TotalFavorited: int64(followed.TotalFavorited),
			WorkCount:      int64(followed.WorkCount),
			FavoriteCount:  int64(followed.FavoriteCount),
		})
	}
	return &relation.RelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}, nil
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	followers, err := model.GetFollowerById(req.UserId)
	if err != nil {
		return &relation.RelationFollowerListResponse{StatusCode: -1, StatusMsg: "获取粉丝列表ID失败"}, err
	}

	followersUser, err := model.GetUsersByIDs(ctx, followers)
	if err != nil {
		return &relation.RelationFollowerListResponse{StatusCode: -1, StatusMsg: "获取粉丝列表USER失败"}, err
	}
	var userList []*user.User
	for _, follower := range followersUser {
		userList = append(userList, &user.User{
			Id:            int64(follower.ID),
			Name:          follower.UserName,
			FollowCount:   int64(follower.FollowingCount),
			FollowerCount: int64(follower.FollowerCount),
			IsFollow:      true,
			// Avatar:          avatar,//to do
			// BackgroundImage: backgroundUrl,//to do
			Signature:      follower.Signature,
			TotalFavorited: int64(follower.TotalFavorited),
			WorkCount:      int64(follower.WorkCount),
			FavoriteCount:  int64(follower.FavoriteCount),
		})
	}
	return &relation.RelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}, nil
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	friends, err := model.GetFriendListById(req.UserId)
	if err != nil {
		return &relation.RelationFriendListResponse{StatusCode: -1, StatusMsg: "获取好友列表ID失败"}, err
	}

	friendsUser, err := model.GetUsersByIDs(ctx, friends)
	if err != nil {
		return &relation.RelationFriendListResponse{StatusCode: -1, StatusMsg: "获取好友列表USER失败"}, err
	}
	var userList []*relation.FriendUser

	fmt.Printf("friendsUser: %v\n", friendsUser)
	for _, friend := range friendsUser {
		var msgType int64
		var msgContent string
		msg, err := model.QueryLastMessageById(req.UserId, int64(friend.ID))
		if err != nil {
			// res := &relation.RelationFriendListResponse{StatusCode: -1, StatusMsg: "获取好友列表MESSAGE失败"}
			// return res, err
			msgContent = ""
			msgType = -1
		} else {
			msgContent = msg.Content
			if msg.FromUserID == req.UserId {
				msgType = 1
			} else {
				msgType = 0
			}
		}
		userList = append(userList, &relation.FriendUser{
			Id:            int64(friend.ID),
			Name:          friend.UserName,
			FollowCount:   int64(friend.FollowingCount),
			FollowerCount: int64(friend.FollowerCount),
			IsFollow:      true,
			// Avatar:          avatar,//to do
			// BackgroundImage: backgroundUrl,//to do
			Signature:      friend.Signature,
			TotalFavorited: int64(friend.TotalFavorited),
			WorkCount:      int64(friend.WorkCount),
			FavoriteCount:  int64(friend.FavoriteCount),
			Message:        msgContent,
			MsgType:        msgType,
		})
	}
	return &relation.RelationFriendListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   userList,
	}, nil
}
