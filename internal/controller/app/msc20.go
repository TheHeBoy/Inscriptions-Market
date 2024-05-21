package app

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/service"
	"gohub/pkg/response"
)

type Msc20Controller struct {
}

var msc20Service = service.Msc20

func (mc *Msc20Controller) GetMsc20ByAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		response.ErrorStr(c, "地址不能为空")
		return
	}
	response.SuccessData(c, msc20Service.GetMsc20ByAddress(address))
}
