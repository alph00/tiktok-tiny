package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// TODO  插入与更新操作： 视频发布相关操作
func SaveVideo(ctx context.Context, video *Video) error {
	err := DB.Clauses(dbresolver.Write).Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//TODO 视频信息写入到video视频表
		insertErr := tx.Create(video).Error
		if insertErr != nil {
			return insertErr
		}
		//TODO 更新 User表中的信息
		updateRes := tx.Debug().Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
		if updateRes.Error != nil {
			return updateRes.Error
		}
		//TODO 数据更新操作不为1
		if updateRes.RowsAffected != 1 {
			return gorm.ErrInvalidValueOfLength
		}
		return nil
	})

	return err
}

// TODO 查询操作： 根据用户id 获取所有的视频信息
func GetVideoListByUserID(ctx context.Context, authorId int64) ([]*Video, error) {
	var pubList []*Video
	err := DB.Clauses(dbresolver.Read).Debug().WithContext(ctx).Model(&Video{}).Where("author_id", uint(authorId)).Find(&pubList).Error
	if err != nil {
		return nil, err
	}
	return pubList, nil
}

// TODO 删除操作： 根据视频id和作者id删除数据库记录
func DelVideoByID(ctx context.Context, videoID int64, authorID int64) error {
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//TODO 根据主键 video_id 删除
		delErr := tx.Unscoped().Delete(&Video{}, videoID).Error
		if delErr != nil {
			return delErr
		}
		// TODO 更新 user表中的作品数量
		updateRes := tx.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count - ?", 1))
		if updateRes.Error != nil {
			return updateRes.Error
		}
		if updateRes.RowsAffected != 1 {
			return gorm.ErrInvalidValueOfLength
		}
		return nil
	})
	return err
}
