package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 通过gin上下文获取用户id
func getId(c *gin.Context) (id int64, ok bool) {
	_id, ok := c.Get("id")
	if !ok {
		return
	}

	id, ok = _id.(int64)
	if !ok {
		log.Println("出现无法解析成64位整数的token")
		return
	}

	return
}
