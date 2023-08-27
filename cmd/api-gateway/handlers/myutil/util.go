package myutil

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type RespMsg struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func ResponseMsg(HttpStatusMsg int, c *app.RequestContext, msg interface{}) {
	c.JSON(HttpStatusMsg, msg)
}
