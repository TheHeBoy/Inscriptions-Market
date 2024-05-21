package app

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/app"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type OrderController struct {
}

var orderService = service.Order

func (lc *OrderController) CreateOrder(c *gin.Context) {
	req := app.CreateOrderReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	err := orderService.Create(req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "创建订单失败")
	} else {
		response.Success(c)
	}
}

func (lc *OrderController) SignOrder(c *gin.Context) {
	req := app.SignOrderReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	err := orderService.SignOrder(req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "签名订单失败")
	} else {
		response.Success(c)
	}
}

func (lc *OrderController) GetListingOrderByTick(c *gin.Context) {
	req := app.GetListingOrderByTickReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}
	resp := orderService.GetListingOrderByTick(req.Tick)
	response.SuccessData(c, resp)

}

func (lc *OrderController) PageBySeller(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		response.ErrorStr(c, "地址不能为空")
		return
	}
	response.SuccessData(c, orderService.GetBySeller(address))
}
