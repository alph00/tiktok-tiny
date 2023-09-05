package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// TODO  用户评论数据模型
type Comment struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_date"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Video      Video          `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID    uint           `gorm:"index:idx_videoid;not null" json:"video_id"`
	User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID     uint           `gorm:"index:idx_userid;not null" json:"user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
	LikeCount  uint           `gorm:"column:like_count;default:0;not null" json:"like_count,omitempty"`
	TeaseCount uint           `gorm:"column:tease_count;default:0;not null" json:"tease_count,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

// TODO  新增一条评论数据，并对所属视频的评论数+1
func InsertComment(ctx context.Context, comment *Comment) error {
	err := DB.Clauses(dbresolver.Write).Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// TODO 由于涉及到两张表comment表以及 video表，引入事务操作 使用 'tx'
		// TODO 1.新增评论 insert
		if tx.Create(comment).Error != nil {
			return tx.Create(comment).Error
		}
		//TODO  2. Video表评论数 +1 update
		resUpdate := tx.Model(&Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if resUpdate.Error != nil {
			return resUpdate.Error
		}
		// TODO 更新的数据条数不是1
		if resUpdate.RowsAffected != 1 {
			return gorm.ErrInvalidValueOfLength
		}

		return nil
	})
	return err
}

// TODO 删除一条评论数据，视频的评论数-1
func DelCommentByID(ctx context.Context, commentID int64, vid int64) error {
	err := DB.Clauses(dbresolver.Write).Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// TODO 由于涉及到两张表comment表以及 video表，引入事务操作 使用 'tx'
		comment := new(Comment)
		//TODO 根据评论id 获取评论
		getErr := tx.First(&comment, commentID).Error
		if getErr != nil {
			return getErr
		}

		//TODO 1. 删除comment表的评论数据,使用的实际上是软删除
		delErr := tx.Where("id = ?", commentID).Delete(&Comment{}).Error
		if delErr != nil {
			return delErr
		}

		//TODO  2.更新video表中的comment_count
		updateRes := tx.Model(&Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if updateRes.Error != nil {
			return updateRes.Error
		}
		// TODO 影响的数据条数不是1
		if updateRes.RowsAffected != 1 {

			return gorm.ErrInvalidValueOfLength
		}

		return nil
	})
	return err
}

// TODO  根据视频id获取指定视频的全部评论内容
func SelectVideoCommentListByVideoID(ctx context.Context, videoID int64) ([]*Comment, error) {
	var comments []*Comment
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Model(&Comment{}).Where("id=", uint(videoID)).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// TODO 根据评论ID获取指定评论的内容
func SelectCommentByCommentID(ctx context.Context, commentID int64) (*Comment, error) {
	comment := new(Comment)
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Where("id = ?", commentID).First(&comment).Error

	if err == nil {
		return comment, nil

	} else if err == gorm.ErrRecordNotFound {

		return nil, nil

	} else {
		return nil, err
	}
}
