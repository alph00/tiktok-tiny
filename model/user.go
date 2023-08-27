package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type User struct {
	gorm.Model
	UserName string `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password string `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	// FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	FollowingCount  uint   `gorm:"default:0;not null" json:"follow_count,omitempty"`                                                           // 关注总数
	FollowerCount   uint   `gorm:"default:0;not null" json:"follower_count,omitempty"`                                                         // 粉丝总数
	Avatar          string `gorm:"type:varchar(256)" json:"avatar,omitempty"`                                                                  // 用户头像
	BackgroundImage string `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"` // 用户个人页顶部大图
	WorkCount       uint   `gorm:"default:0;not null" json:"work_count,omitempty"`                                                             // 作品数
	FavoriteCount   uint   `gorm:"default:0;not null" json:"favorite_count,omitempty"`                                                         // 喜欢数
	TotalFavorited  uint   `gorm:"default:0;not null" json:"total_favorited,omitempty"`                                                        // 获赞总量
	Signature       string `gorm:"type:varchar(256)" json:"signature,omitempty"`                                                               // 个人简介
}

func (User) TableName() string {
	return "users"
}

func GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetUserByID(ctx context.Context, userID int64) (*User, error) {
	res := new(User)
	if err := DB.Clauses(dbresolver.Read).WithContext(ctx).First(&res, userID).Error; err == nil {
		return res, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func CreateUsers(ctx context.Context, users []*User) error {
	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(users).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func CreateUser(ctx context.Context, user *User) error {
	err := DB.Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetUserByName(ctx context.Context, userName string) (*User, error) {
	res := new(User)
	if err := DB.Clauses(dbresolver.Read).WithContext(ctx).Select("id, user_name, password").Where("user_name = ?", userName).First(&res).Error; err == nil {
		return res, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func GetPasswordByUsername(ctx context.Context, userName string) (*User, error) {
	user := new(User)
	if err := DB.Clauses(dbresolver.Read).WithContext(ctx).
		Select("password").Where("user_name = ?", userName).
		First(&user).Error; err == nil {
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
