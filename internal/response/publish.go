package response

import (
	"github.com/alph00/tiktok-tiny/kitex_gen/publish"
)

type PublishAction struct {
	Base
}

type PublishList struct {
	Base
	VideoList []*publish.Video `json:"video_list"`
}

type Feed struct {
	Base
	NextTime  int64            `json:"next_time"`
	VideoList []*publish.Video `json:"video_list"`
}
