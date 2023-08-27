package model

import (
	"context"
	"fmt"
	"github.com/alph00/tiktok-tiny/pkg/errno"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// TODO  插入与更新操作： 视频发布相关操作
func CreateVideo(ctx context.Context, video *Video) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//TODO 视频信息写入到video视频表
		err := tx.Create(video).Error
		fmt.Println(err)
		if err != nil {
			return err
		}
		//TODO 更新 User表中的信息
		res := tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})

	return err
}

// TODO 查询操作： 根据用户id 获取所有的视频信息
func GetVideoListByUserID(ctx context.Context, authorId int64) ([]*Video, error) {
	var pubList []*Video
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Video{}).Where(&Video{AuthorID: uint(authorId)}).Find(&pubList).Error
	if err != nil {
		return nil, err
	}
	return pubList, nil
}

// TODO 删除操作： 根据视频id和作者id删除数据库记录
func DelVideoByID(ctx context.Context, videoID int64, authorID int64) error {
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//TODO 根据主键 video_id 删除
		err := tx.Unscoped().Delete(&Video{}, videoID).Error
		if err != nil {
			return err
		}
		// TODO 更新 user表中的作品数量
		res := tx.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})
	return err
}
