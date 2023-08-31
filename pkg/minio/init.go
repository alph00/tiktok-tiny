package minio

import (
	"log"
	"time"

	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
	expire      time.Duration
)

func init() {
	minioConfig := viper.Read("minio")
	expire = time.Second * minioConfig.GetDuration("expire")
	var err error
	minioClient, err = minio.New(minioConfig.GetString("endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.GetString("accessKeyID"), minioConfig.GetString("secretAccessKey"), ""),
		Secure: minioConfig.GetBool("useSSL")})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	}
	log.Printf("%#v\n", minioClient)
}
