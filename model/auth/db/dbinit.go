package db

import (
	"fmt"
	"time"

	"github.com/alph00/tiktok-tiny/package/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_db    *gorm.DB
	config = viper.Init("db")
	// zapLogger = zap.InitLogger()
)

func getDsn(driverWithRole string) string {
	username := config.Viper.GetString(fmt.Sprintf("%s.username", driverWithRole))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", driverWithRole))
	host := config.Viper.GetString(fmt.Sprintf("%s.host", driverWithRole))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", driverWithRole))
	Dbname := config.Viper.GetString(fmt.Sprintf("%s.database", driverWithRole))

	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	return dsn
}
func dbinit() {
	// TODO:日志打印

	dsn1 := getDsn("mysql.source")
	var err error

	_db, err := gorm.Open(mysql.Open(dsn1), &gorm.Config{
		// TODO:日志打印
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")

	}

	// AutoMigrate会创建表，缺失的外键，约束，列和索引。如果大小，精度，是否为空，可以更改，则AutoMigrate会改变列的类型。出于保护您数据的目的，它不会删除未使用的列
	// 刷新数据库的表格，使其保持最新。即如果我在旧表的基础上增加一个字段age，那么调用autoMigrate后，旧表会自动多出一列age，值为空
	if err := _db.AutoMigrate(&Auth{}); err != nil {
		panic(err.Error())
	}

	db, err := _db.DB()
	if err != nil {
		panic(err.Error())
		// zapLogger.Fatalln(err.Error())
	}
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
}
func GetDB() *gorm.DB {
	return _db
}
