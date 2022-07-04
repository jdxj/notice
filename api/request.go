package api

import (
	"github.com/gin-gonic/gin"

	"github.com/jdxj/notice/logger"
)

type Request struct {
	Data interface{} `json:"data"`
}

func ParseRequest(c *gin.Context, data interface{}) (*Request, error) {
	r := &Request{
		Data: data,
	}

	err := c.ShouldBind(r)
	if err != nil {
		logger.Errorf("bind request err: %s", err)

	}
}
