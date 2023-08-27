package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

func Read(configName string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName(configName)
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	if err := v.ReadInConfig(); err != nil {
		// log.Fatalf("errno is %+v", err)
		fmt.Printf("err: %v\n", err)
	}
	return v
}
