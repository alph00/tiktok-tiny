package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Video struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null;index:idx_create" json:"created_at,omitempty"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	AuthorID      uint           `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string         `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string         `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount uint           `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint           `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string         `gorm:"type:varchar(50);not null" json:"title,omitempty"`
	// Author        User   `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

// TODO 获取最近发表的视频
func GetLatestTimeVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	//TODO 数据库查查询列表信息
	results := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}
	//TODO 查询语句
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Limit(limit).Order("created_at desc").Where("created_at < ?", time.UnixMilli(*latestTime)).Find(&results).Error
	//TODO 异常处理
	if err != nil {
		return nil, err
	}
	return results, nil
}

// TODO 根据视频id获取视频信息
func GetVideoInfoById(ctx context.Context, videoId int64) (*Video, error) {
	//TODO 数据库实例
	video := new(Video)
	//TODO 查询语句
	err := DB.Clauses(dbresolver.Read).WithContext(ctx).Where("id", videoId).First(&video).Error
	//TODO 异常处理
	if err == nil {
		return video, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
