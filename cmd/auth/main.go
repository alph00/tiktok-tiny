package main

import (
	"log"

	auth "github.com/alph00/tiktok-tiny/kitex_gen/auth/authservice"
	"github.com/alph00/tiktok-tiny/package/viper"
)

var (
	config = viper.Init("user")
	// serviceName = config.Viper.GetString("server.name")
	// serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	// etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey = config.Viper.GetString("JWT.signingKey")
	// logger      = zap.InitLogger()
)

func init() {
	Init(signingKey)
}

func main() {
	svr := auth.NewServer(new(AuthServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
