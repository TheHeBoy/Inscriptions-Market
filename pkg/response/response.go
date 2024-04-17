// Package response 响应处理工具
package response

import (
	"gohub/pkg/errorcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonResult struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResult{
		Code: 200,
		Data: nil,
		Msg:  "",
	})
}

func SuccessData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, CommonResult{
		Code: 200,
		Data: data,
		Msg:  "",
	})
}

func Error(c *gin.Context, errorCode errorcode.ErrorCode) {
	c.JSON(http.StatusOK, CommonResult{
		Code: errorCode.Code,
		Data: nil,
		Msg:  errorCode.Msg,
	})
}

func ErrorCustom(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, CommonResult{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
