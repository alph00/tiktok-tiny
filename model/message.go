package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         int64     `gorm:"primarykey"`
	FromUserID int64     `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUserID   int64     `gorm:"index:idx_userid_from;index:idx_userid_to;not null" json:"to_user_id"`
	Content    string    `gorm:"type:varchar(255);not null" json:"content"`
	CreatedAt  time.Time `gorm:"index;not null" json:"created_at"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Message) TableName() string {
	return "messages"
}

func CreateMessage(message *Message) error {
	if result := DB.Create(&message); result.Error != nil {
		return result.Error
	}
	return nil
}

func QueryMessageList(date *string, fromUserId int64, ToUserId int64) ([]*Message, error) {
	// fmt.Println(*date)
	var MessageList []*Message
	result := DB.Where("( (from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?) ) and created_at > ?", fromUserId, ToUserId, ToUserId, fromUserId, date).Order("created_at asc").Find(&MessageList)
	if result.Error != nil {
		return nil, result.Error
	}
	// fmt.Println(MessageList)
	return MessageList, nil
}

func QueryLastMessageById(fromUserId int64, ToUserId int64) (*Message, error) {
	// fmt.Println(*date)
	var MessageLast *Message
	result := DB.Where("( (from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?) ) ", fromUserId, ToUserId, ToUserId, fromUserId).Order("created_at desc").Take(&MessageLast)
	if result.Error != nil {
		return nil, result.Error
	}
	// fmt.Println(MessageList)
	return MessageLast, nil
}
