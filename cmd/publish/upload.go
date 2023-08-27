package main

import (
	"bytes"
	"fmt"
	tool "github.com/alph00/tiktok-tiny/internal/tools"
	"github.com/alph00/tiktok-tiny/pkg/minio"
	"github.com/bytedance/gopkg/util/logger"
)

// TODO 上传视频至 minio
func uploadVideo(data []byte, title string) (string, error) {

	//TODO 读取视频的文件流
	filereader := bytes.NewReader(data)
	//TODO 文件的类型
	contentType := "application/mp4"
	//TODO minio 文件上传
	uploadSize, err := minio.UploadFileByIO(minio.VideoBucketName, title, filereader, int64(len(data)), contentType)
	if err != nil {
		logger.Errorf("视频上传minio失败：%v", err.Error())
		return "", err
	}
	logger.Infof("视频文件大小为：%v", uploadSize)
	//TODO 获取minio中的文件路径
	playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, title)
	if err != nil {
		logger.Errorf("服务内部异常：视频获取失败：%s", err.Error())
		return "", err
	}
	logger.Infof("上传视频路径：%v", playUrl)
	return playUrl, nil
}

// TODO 上传封面至 minio
func uploadCover(playUrl string, title string) error {
	//TODO 截取视频的第一帧为视频的封面
	imgBuffer, err := tool.GetSnapshotImageBuffer(playUrl, 1)
	fmt.Println(err)
	if err != nil {
		logger.Errorf("服务内部异常：封面获取失败：%s", err.Error())
		return err
	}
	//TODO 封面图片字节流写入
	var imgByte []byte
	imgBuffer.Write(imgByte)
	contentType := "image/png"
	//TODO 封面上传至minio
	uploadSize, err := minio.UploadFileByIO(minio.CoverBucketName, title, imgBuffer, int64(imgBuffer.Len()), contentType)
	if err != nil {
		logger.Errorf("封面上传至minio失败：%v", err.Error())
		return err
	}
	//TODO 获取上传封面的URL路径
	coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, title)
	if err != nil {
		logger.Errorf("封面获取链接失败：%v", err.Error())
		return err
	}
	logger.Infof("上传封面URL路径：%v", coverUrl)
	logger.Infof("封面文件大小为：%v", uploadSize)

	return nil
}

// TODO  上传视频并获取视频封面
func VideoPublish(data []byte, videoTitle string, coverTitle string) error {
	//TODO 上传视频
	playUrl, err := uploadVideo(data, videoTitle)
	if err != nil {
		return err
	}
	//TODO 根据文件路径获取视频封面 并上传到minio
	err = uploadCover(playUrl, coverTitle)
	if err != nil {
		//TODO 视频上传失败的异常
		logger.Error(err)
		return err
	}
	return nil
}
