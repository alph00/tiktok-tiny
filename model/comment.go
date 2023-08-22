package model

import (
	"context"
	"github.com/alph00/tiktok-tiny/pkg/errno"
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
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// TODO 由于涉及到两张表comment表以及 video表，引入事务操作 使用 'tx'
		// TODO 1.新增评论 insert
		err := tx.Create(comment).Error
		if err != nil {
			return err
		}
		//TODO  2. Video表评论数 +1 update
		res := tx.Model(&Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// TODO 删除一条评论数据，视频的评论数-1
func DelCommentByID(ctx context.Context, commentID int64, vid int64) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// TODO 由于涉及到两张表comment表以及 video表，引入事务操作 使用 'tx'
		comment := new(Comment)
		if err := tx.First(&comment, commentID).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		//TODO 1. 删除comment表的评论数据,使用的实际上是软删除
		err := tx.Where("id = ?", commentID).Delete(&Comment{}).Error
		if err != nil {
			return err
		}

		//TODO  2.更新video表中的comment_count
		res := tx.Model(&Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// TODO  根据视频id获取指定视频的全部评论内容
func SelectVideoCommentListByVideoID(ctx context.Context, videoID int64) ([]*Comment, error) {
	var comments []*Comment
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Comment{}).Where(&Comment{VideoID: uint(videoID)}).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// TODO 根据评论ID获取指定评论的内容
func SelectCommentByCommentID(ctx context.Context, commentID int64) (*Comment, error) {
	comment := new(Comment)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id = ?", commentID).First(&comment).Error; err == nil {
		return comment, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
