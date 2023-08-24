// 用户登录的数据库逻

package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Auth struct {
	gorm.Model
	UserName string `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password string `gorm:"type:varchar(256);not null" json:"password,omitempty"`
}

func (Auth) TableName() string {
	return "auths"
}

// CreateUser
//
//	@Description: 新增一条用户数据
//	@Date 2023-02-22 11:46:43
//	@param ctx 数据库操作上下文
//	@param user 用户数据
//	@return error
func CreateUser(ctx context.Context, user *Auth) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetUserByName
//
//	@Description: 根据用户名获取用户数据列表
//	@Date 2023-01-21 17:15:17
//	@param ctx 数据库操作上下文
//	@param userName 用户名
//	@return []*User 用户数据列表
//	@return error
func GetUserByName(ctx context.Context, userName string) (*Auth, error) {
	res := new(Auth)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("id,password").Where("user_name = ?", userName).First(&res).Error; err == nil {
		return res, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
