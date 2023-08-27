package viper

import (
	"log"

	V "github.com/spf13/viper"
)

// TODO 配置项维护与管理
type Config struct {
	Viper *V.Viper
}

// TODO 读取项目配置文件初始化工作
func Init(configName string) Config {
	config := Config{Viper: V.New()}
	v := config.Viper
	//TODO 配置文件格式yml
	v.SetConfigType("yml")
	//TODO 配置文件名字
	v.SetConfigName(configName)
	//TODO 配置文件的路径
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	//TODO 读取相关的配置
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("errno is %+v", err)
	}
	return config
}
