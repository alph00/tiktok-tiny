package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alph00/tiktok-tiny/model"
)

func consume() error {
	msgs, err := RelationMq.ConsumeSimple()
	if err != nil {
		fmt.Println(err.Error())
		log.Printf("RelationMQ Err: %s", err.Error())
	}
	// 将消息队列的消息全部取出
	for msg := range msgs {
		fc := new(model.Relation)
		// 解析json
		if err = json.Unmarshal(msg.Body, &fc); err != nil {
			fmt.Println("json unmarshal error:" + err.Error())
		}
		fmt.Printf("==> Get new message: %v\n", fc)
		// 将结构体存入mysql
		// if model.IsFollow(fc.FollowedId, fc.FollowerId) {
		// 	return &relation.RelationActionResponse{StatusCode: 0, StatusMsg: "已经关注过该用户"}, nil //这样算成功吗？
		// }
		err := model.Follow(fc.FollowedId, fc.FollowerId)
		if err != nil {
			fmt.Println("关注信息存入数据库发生错误:" + err.Error())
			return err
		}
		model.AddFollowerCount(fc.FollowedId, 1)
		model.AddFollowingCount(fc.FollowerId, 1)

		// if err = model.CreateVideoFavorite(context.Background(), fc); err != nil {
		// 	logger.Errorf("json unmarshal error: %s", err.Error())
		// 	fmt.Println("点赞结构体存入数据库发生错误:" + err.Error())
		// 	continue
		// }
		//err = redis.UnlockByMutex(context.Background(), redis.FavoriteMutex)
		//if err != nil {
		//	logger.Errorf("Redis mutex unlock error: %s", err.Error())
		//	return err
		//}
		if !autoAck {
			err := msg.Ack(true)
			if err != nil {
				// logger.Errorf("ack error: %s", err.Error())
				return err
			}
		}
	}
	return nil

}
func Delconsume() error {
	msgs, err := RelationDelMq.ConsumeSimple()
	if err != nil {
		fmt.Println(err.Error())
		log.Printf("RelationMQ Err: %s", err.Error())
	}
	// 将消息队列的消息全部取出
	for msg := range msgs {
		fc := new(model.Relation)
		// 解析json
		if err = json.Unmarshal(msg.Body, &fc); err != nil {
			fmt.Println("json unmarshal error:" + err.Error())
		}
		fmt.Printf("==> Get new message: %v\n", fc)
		// 将结构体存入mysql
		// if model.IsFollow(fc.FollowedId, fc.FollowerId) {
		// 	return &relation.RelationActionResponse{StatusCode: 0, StatusMsg: "已经关注过该用户"}, nil //这样算成功吗？
		// }
		err := model.UnFollow(fc.FollowedId, fc.FollowerId)
		if err != nil {
			fmt.Println("关注信息删除数据库发生错误:" + err.Error())
			return err
		}
		model.ReduceFollowerCount(fc.FollowedId, 1)
		model.ReduceFollowingCount(fc.FollowerId, 1)

		// if err = model.CreateVideoFavorite(context.Background(), fc); err != nil {
		// 	logger.Errorf("json unmarshal error: %s", err.Error())
		// 	fmt.Println("点赞结构体存入数据库发生错误:" + err.Error())
		// 	continue
		// }
		//err = redis.UnlockByMutex(context.Background(), redis.FavoriteMutex)
		//if err != nil {
		//	logger.Errorf("Redis mutex unlock error: %s", err.Error())
		//	return err
		//}
		if !autoAck {
			err := msg.Ack(true)
			if err != nil {
				// logger.Errorf("ack error: %s", err.Error())
				return err
			}
		}
	}
	return nil

}
