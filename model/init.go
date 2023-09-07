package model

import (
	"github.com/alph00/tiktok-tiny/pkg/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	mysqlConfig := viper.Read("database")
	dsn := mysqlConfig.GetString("mysql.username") + ":" + mysqlConfig.GetString("mysql.password") + "@tcp(" + mysqlConfig.GetString("mysql.host") + ":" + mysqlConfig.GetString("mysql.port") + ")/" + mysqlConfig.GetString("mysql.database") + "?charset=utf8&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
}
