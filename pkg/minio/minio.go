package minio

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

func CreateBucket(bucketName string) error {
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Println("创建bucket错误: ", err)
		exists, errEx := minioClient.BucketExists(ctx, bucketName)
		if exists && errEx != nil {
			log.Printf("bucket: %s已经存在", bucketName)
		} else {
			return errEx
		}

	}
	return nil
}
func listBucket() {
	ctx := context.Background()
	buckets, _ := minioClient.ListBuckets(ctx)
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
}

func UploadFileByIO(bucketName string, objectName string, reader io.Reader, size int64, contentType string) (int64, error) {
	if uploadInfo, err := minioClient.PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return -1, err
	} else {
		return uploadInfo.Size, nil
	}
}

func GetFileTemporaryURL(bucketName, objectName string) (string, error) {
	if presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expire, nil); err != nil {
		return "", err
	} else {
		return presignedURL.String(), nil
	}
}
