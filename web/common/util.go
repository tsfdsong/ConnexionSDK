package common

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	C *gin.Context
}

func (c *Context) Response(httpCode int, code int, data interface{}) {
	c.C.JSON(httpCode, gin.H{
		"code": code,
		"msg":  ErrorMap[code],
		"data": data,
	})
	return
}
