package service

import (
	"context"
	"time"

	message "github.com/alph00/tiktok-tiny/kitex_gen/message"
	"github.com/alph00/tiktok-tiny/model"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageChat implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageChat(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	userId := req.Id
	// fmt.Printf("userId: %v\n", userId)
	toUId := req.ToUserId

	lastTime := req.PreMsgTime
	format := "2006-01-02 15:04:05"
	t := time.Unix(lastTime/1000, 0)
	// fmt.Printf("lastTime: %v\n", lastTime)
	searchTime := t.Format(format)

	var messageList []*model.Message
	messageList, err = model.QueryMessageList(&searchTime, userId, toUId)
	if err != nil {
		return nil, err
	}

	var res []*message.Message
	for _, msg := range messageList {
		chatMsg := &message.Message{
			Id:         msg.ID,
			Content:    msg.Content,
			FromUserId: msg.FromUserID,
			ToUserId:   msg.ToUserID,
			CreateTime: msg.CreatedAt.UnixMilli(),
		}
		res = append(res, chatMsg)
	}
	return &message.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: res,
	}, nil
}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	userId := req.Id

	toUId, actType := req.ToUserId, req.ActionType

	if userId == toUId {
		res := &message.MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "消息发送失败：不能给自己发送消息",
		}
		return res, nil
	}

	// relation, err := model.GetRelationByUserIDs(ctx, userId, toUId)
	// if relation == nil {
	// 	logger.Errorf("消息发送失败：非朋友关系，无法发送")
	// 	res := &message.MessageActionResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "消息发送失败：非朋友关系，无法发送",
	// 	}
	// 	return res, nil
	// }

	//to do 加密 编码
	if actType == 1 {
		err := model.CreateMessage(&model.Message{
			FromUserID: userId,
			ToUserID:   toUId,
			Content:    req.Content,
		})
		if err != nil {
			res := &message.MessageActionResponse{
				StatusCode: -1,
				StatusMsg:  "failure",
			}
			return res, nil
		}
	} else {
		res := &message.MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "failure",
		}
		return res, nil
	}
	res := &message.MessageActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return res, nil
}
