package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alph00/tiktok-tiny/model"
	"github.com/bytedance/gopkg/util/logger"
)

func consume() error {
	msgs, err := FavoriteMq.ConsumeSimple()
	if err != nil {
		fmt.Println(err.Error())
		log.Printf("FavoriteMQ Err: %s", err.Error())
	}
	// 将消息队列的消息全部取出
	for msg := range msgs {
		fc := new(model.FavoriteVideoRelation)
		// 解析json
		if err = json.Unmarshal(msg.Body, &fc); err != nil {
			fmt.Println("json unmarshal error:" + err.Error())
		}
		fmt.Printf("==> Get new message: %v\n", fc)
		// 将结构体存入redis
		if err = model.CreateVideoFavorite(context.Background(), fc); err != nil {
			logger.Errorf("json unmarshal error: %s", err.Error())
			fmt.Println("点赞结构体存入数据库发生错误:" + err.Error())
			continue
		}
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
