package minio

import (
	"context"
	"log"
	"time"

	"github.com/alph00/tiktok-tiny/pkg/viper"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient     *minio.Client
	expire          time.Duration
	Video           string
	Cover           string
	Avatar          string
	BackgroundImage string
)

func init() {
	minioConfig := viper.Read("minio")
	expire = time.Second * minioConfig.GetDuration("expire")
	Video = minioConfig.GetString("videoBucketName")
	Cover = minioConfig.GetString("coverBucketName")
	Avatar = minioConfig.GetString("avatarBucketName")
	BackgroundImage = minioConfig.GetString("backgroundImageBucketName")
	var err error
	minioClient, err = minio.New(minioConfig.GetString("endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.GetString("accessKeyID"), minioConfig.GetString("secretAccessKey"), ""),
		Secure: minioConfig.GetBool("useSSL")})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	}

	exist, _ := minioClient.BucketExists(context.TODO(), Video)
	if !exist {
		if err := CreateBucket(Video); err != nil {
			panic(err)
		}
	}

	exist, _ = minioClient.BucketExists(context.TODO(), Cover)
	if !exist {
		if err := CreateBucket(Cover); err != nil {
			panic(err)
		}
	}
	exist, _ = minioClient.BucketExists(context.TODO(), Avatar)
	if !exist {

		if err := CreateBucket(Avatar); err != nil {
			panic(err)
		}
	}

	exist, _ = minioClient.BucketExists(context.TODO(), BackgroundImage)
	if !exist {
		if err := CreateBucket(BackgroundImage); err != nil {
			panic(err)
		}
	}
}
