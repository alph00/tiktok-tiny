package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type FavoriteVideoRelation struct {
	// Video   Video `gorm:"foreignkey:VideoID;" json:"video,omitempty"`
	UserID  uint `gorm:"index:idx_userid;not null" json:"user_id"`
	VideoID uint `gorm:"index:idx_videoid;not null" json:"video_id"`
	// User    User  `gorm:"foreignkey:UserID;" json:"user,omitempty"`

}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func CreateVideoFavorite(ctx context.Context, favorite_video_data *FavoriteVideoRelation) error {
	fmt.Println("开始访问数据库，创建一条点赞数据")
	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(favorite_video_data).Error
		if err != nil {
			return err
		}
		fmt.Println("创建一条点赞数据成功")
		res := tx.Model(&Video{}).Where("id = ?", favorite_video_data.VideoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}
		res2 := tx.Model(&User{}).Where("id = ?", favorite_video_data.UserID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res2.Error != nil {
			return err
		}
		fmt.Println("修改点赞人user表成功，创建一条点赞数据成功")

		//4.改变 user 表中的 total_favorited  TODO:这里涉及到查询authorID，基本框架没问题
		// var video Feed
		// DB.Find(&video,favorite_video_data.VideoID)
		// authorID :=video.
		// res3 := tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		// if res3.Error != nil {
		// 	return err
		// }
		return nil
	})
	return err
}

func IsFavorite(uId int64, vId int64) bool {
	var count int64
	DB.Model(&FavoriteVideoRelation{}).Where("user_id = ? and video_id = ? ", uId, vId).Count(&count)
	if count == 0 {
		return false
	} else {
		return true
	}
}

func DelVideoFavorite(ctx context.Context, favorite_video_data *FavoriteVideoRelation) error {
	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		FavoriteVideoRelation := new(FavoriteVideoRelation)
		if err := tx.Where("user_id = ? and video_id = ?", favorite_video_data.UserID, favorite_video_data.VideoID).First(&FavoriteVideoRelation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}
		err := tx.Unscoped().Where("user_id = ? and video_id = ?", favorite_video_data.UserID, favorite_video_data.VideoID).Delete(&FavoriteVideoRelation).Error

		if err != nil {
			return err
		}
		fmt.Println("创建一条点赞数据成功")

		res := tx.Model(&Video{}).Where("id = ?", favorite_video_data.VideoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		res2 := tx.Model(&User{}).Where("id = ?", favorite_video_data.UserID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res2.Error != nil {
			return err
		}
		fmt.Println("修改点赞人user表成功，创建一条点赞数据成功")

		//4.改变 user 表中的 total_favorited  TODO:这里涉及到查询authorID，基本框架没问题
		// var video Feed
		// DB.Find(&video,favorite_video_data.VideoID)
		// authorID :=video.
		// res = tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		// if res.Error != nil {
		// 	return err
		// }
		// if res.RowsAffected != 1 {
		// 	return errno.ErrDatabase
		// }

		return nil
	})
	return err
}

func ShowVideoFavorite(ctx context.Context, UserId uint) ([]uint, error) {
	// var VideoList []*Video
	fmt.Println("进行FavoriteList函数   Favorite表查询 ")
	var VideoIDList []uint
	rows, err := DB.Table("user_favorite_videos").Select("COALESCE(video_id, ?)", -1).Where("user_id=?", UserId).Rows()
	if err != nil {
		// 处理错误
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var video_id uint
		if err := rows.Scan(&video_id); err != nil {
			fmt.Println("扫描错误")
		}
		fmt.Println("VideoID:", video_id)
		VideoIDList = append(VideoIDList, video_id)
	}

	if err := rows.Err(); err != nil {
		// 处理迭代错误
		fmt.Println("迭代错误")
	}

	// result := DB.Model(&FavoriteVideoRelation{}).Where("user_id=?", UserId).Find(&VideoIDList)
	// // result := DB.Select("video_id").Where("user_id=?", UserId).Find(&VideoIDList)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }
	return VideoIDList, nil

	// return nil, err

}
