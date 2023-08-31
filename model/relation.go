package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	FollowedId int64 `gorm:"column:followed_id"`
	FollowerId int64 `gorm:"column:follower_id"`
}

func (r *Relation) TableName() string {
	return "relations"
}

func GetFollowerById(myId int64) ([]int64, error) {
	followers := make([]int64, 0)
	err := DB.Model(&Relation{}).Select("follower_id").Where("followed_id=?", myId).Pluck("follower_id", &followers).Error
	fmt.Printf("followers: %v\n", followers)
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return followers, err
}

func GetFollowedById(myId int64) ([]int64, error) {
	followeds := make([]int64, 0)
	err := DB.Model(&Relation{}).Select("followed_id").Where("follower_id=?", myId).Pluck("followed_id", &followeds).Error
	fmt.Printf("followeds: %v\n", followeds)
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return followeds, err
}

func Follow(followedId, followerId int64) error {
	Relation := Relation{
		FollowerId: followerId,
		FollowedId: followedId,
	}
	return DB.Create(&Relation).Error
}

func UnFollow(followedId, followerId int64) error {
	Relation := Relation{
		FollowerId: followerId,
		FollowedId: followedId,
	}
	return DB.Where("follower_id = ? AND followed_id = ?", followerId, followedId).Delete(&Relation).Error
}

func GetFriendListById(myId int64) ([]int64, error) {
	var followeds []int64
	err := DB.Table("relations a").Select("a.follower_id").Where("a.followed_id = ?", myId).Joins("INNER JOIN relations b ON a.follower_id = b.followed_id AND a.followed_id = b.follower_id").Pluck("followed_id", &followeds).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return followeds, err
}

func IfFriend(myId int64, anothorId int64) bool {
	var count1, count2 int64
	DB.Model(&Relation{}).Where("followed_id = ? and follower_id = ? ", myId, anothorId).Count(&count1)
	DB.Model(&Relation{}).Where("follower_id = ? and followed_id = ? ", myId, anothorId).Count(&count2)
	if count1 == 0 || count2 == 0 {
		return false
	} else {
		return true
	}
}

func IsFollow(followedId, followerId int64) bool {
	var count int64
	DB.Model(&Relation{}).Where("followed_id = ? and follower_id = ? ", followedId, followerId).Count(&count)
	if count == 0 {
		return false
	} else {
		return true
	}
}
