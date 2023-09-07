package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func SaveVideo(ctx context.Context, video *Video) error {
	err := DB.Clauses(dbresolver.Write).Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		insertErr := tx.Create(video).Error
		if insertErr != nil {
			return insertErr
		}
		updateRes := tx.Debug().Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
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

func GetVideoListByUserID(ctx context.Context, authorId int64) ([]*Video, error) {
	var pubList []*Video
	err := DB.Clauses(dbresolver.Read).Debug().WithContext(ctx).Model(&Video{}).Where("author_id", uint(authorId)).Find(&pubList).Error
	if err != nil {
		return nil, err
	}
	return pubList, nil
}

func DelVideoByID(ctx context.Context, videoID int64, authorID int64) error {
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		delErr := tx.Unscoped().Delete(&Video{}, videoID).Error
		if delErr != nil {
			return delErr
		}
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
