package main

import (
	"fmt"
	tool "github.com/alph00/tiktok-tiny/internal/tools"
	"testing"
)

func TestUploadCover(t *testing.T) {
	fmt.Println("test")
	url := "http://127.0.0.1:9000/tiktok-videos/1_%E6%8A%96%E9%9F%B3%E8%A7%86%E9%A2%91_1692688917837.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20230822%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230822T072157Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=00bfde2a7d08f747bdd6c048e690cce95a7d102ea04b594410d0d0be499f56af"
	//err := uploadCover(url, "testvover")
	imgBuffer, err := tool.GetSnapshot(url, "test", 1)
	fmt.Println(err)
	fmt.Println(imgBuffer)
	fmt.Println(err)
}
