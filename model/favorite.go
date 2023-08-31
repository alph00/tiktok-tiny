package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// FavoriteVideoRelation
//
//	@Description: 用户与视频的点赞关系数据模型
type FavoriteVideoRelation struct {
	// Video   Video `gorm:"foreignkey:VideoID;" json:"video,omitempty"`
	UserID  uint `gorm:"index:idx_userid;not null" json:"user_id"`
	VideoID uint `gorm:"index:idx_videoid;not null" json:"video_id"`
	// User    User  `gorm:"foreignkey:UserID;" json:"user,omitempty"`

}
type Fav_vid_API struct {
	VideoID []*uint
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

// CreateVideoFavorite
//
//	@Description: 创建一条用户点赞数据
//	@param ctx 数据库操作上下文
//	@param userID 用户id
//	@param videoID 视频id
//	@return error
func CreateVideoFavorite(ctx context.Context, favorite_video_data *FavoriteVideoRelation) error {
	// TODO:这里使用事务
	fmt.Println("开始访问数据库，创建一条点赞数据")

	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1. 新增点赞数据
		err := tx.Create(favorite_video_data).Error
		if err != nil {
			return err
		}
		fmt.Println("创建一条点赞数据成功")

		//2.改变 feed 表中的 favorite count  TODO：feed
		// res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }

		// if res.RowsAffected != 1 {
		// 	// 影响的数据条数不是1
		// 	return errno.ErrDatabase
		// }

		//3.改变 user 表中的 favorite count
		res := tx.Model(&User{}).Where("id = ?", favorite_video_data.UserID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return err
		}
		// if res.RowsAffected != 1 {
		// 	return errno.ErrDatabase
		// }
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

// DelVideoFavorite
//
//	@Description: 删除一条用户点赞数据
//	@param ctx 数据库操作上下文
//	@param userID 用户id
//	@param videoID 视频id
//	@return error
func DelVideoFavorite(ctx context.Context, favorite_video_data *FavoriteVideoRelation) error {
	// TODO:这里使用事务

	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）

		// 先进行查找
		FavoriteVideoRelation := new(FavoriteVideoRelation)
		if err := tx.Where("user_id = ? and video_id = ?", favorite_video_data.UserID, favorite_video_data.VideoID).First(&FavoriteVideoRelation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}
		//1. 删除点赞数据
		// 因为FavoriteVideoRelation中包含了gorm.Model所以拥有软删除能力
		// 而tx.Unscoped().Delete()将永久删除记录
		err := tx.Unscoped().Where("user_id = ? and video_id = ?", favorite_video_data.UserID, favorite_video_data.VideoID).Delete(&FavoriteVideoRelation).Error

		if err != nil {
			return err
		}
		fmt.Println("创建一条点赞数据成功")

		//2.改变 feed 表中的 favorite count  TODO：feed
		// res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }

		// if res.RowsAffected != 1 {
		// 	// 影响的数据条数不是1
		// 	return errno.ErrDatabase
		// }

		//3.改变 user 表中的 favorite count
		res := tx.Model(&User{}).Where("id = ?", favorite_video_data.UserID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return err
		}
		// if res.RowsAffected != 1 {
		// 	return errno.ErrDatabase
		// }
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

// ShowVideoFavorite
//
//	@Description: 查找用户点赞数据
//	@param ctx 数据库操作上下文
//	@param userID 用户id
//	@return error   Real   添加feed后进行修改
// func ShowVideoFavorite(ctx context.Context, UserId uint) ([]*Video, error) {
// 	var VideoList []*Video
// 	var VideoIDList []*uint
// 	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
// 		//1. 在favotite--表里查找   userid---->videoID    返回VideoIDLists
// 		result := tx.Where("user_id=?", UserId).Find(&VideoIDList)
// 		if result.Error != nil {
// 			return nil, result.Error
// 		}
// 		return VideoIDList, nil
// 		//2.查找 feed 表中的 feeds
// 		// res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
// 		// if res.Error != nil {
// 		// 	return res.Error
// 		// }

// 		return nil
// 	})
// 	return err

// }
func ShowVideoFavorite(ctx context.Context, UserId uint) ([]uint, error) {
	// var VideoList []*Video
	fmt.Println("进行FavoriteList函数   Favorite表查询 ")
	var VideoIDList []uint
	rows, err := DB.Table("user_favorite_videos").Select("COALESCE(video_id, ?)", -1).Where("user_id=?", UserId).Rows()
	if err != nil {
		// 处理错误
		return nil, err
	}

	defer rows.Close() // 一定要记得关闭迭代器

	for rows.Next() {
		var video_id uint
		if err := rows.Scan(&video_id); err != nil {
			// 处理扫描错误
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
