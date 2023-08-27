package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"user_name" column:"user_name"`
	Email    string `json:"email" column:"email"`
	Password string `json:"password" column:"password"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUsers(users []*User) error {
	return DB.Create(users).Error
}

func FindUserByNameOrEmail(userName, email string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.Where(DB.Or("user_name = ?", userName).
		Or("email = ?", email)).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func CheckUser(account, password string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.Where(DB.Or("user_name = ?", account).
		Or("email = ?", account)).Where("password = ?", password).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
