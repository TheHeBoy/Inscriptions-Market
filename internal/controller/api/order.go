package api

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/api"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type OrderController struct {
}

var orderService = service.Order

func (lc *OrderController) CreateOrder(c *gin.Context) {
	req := api.CreateOrderReq{}
	if ok := validators.Validate(c, &req, api.CreateOrderVal); !ok {
		return
	}

	err := orderService.Create(req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "创建失败")
	} else {
		response.Success(c)
	}
}

func (lc *OrderController) PageListingOrder(c *gin.Context) {
	req := api.PageListingOrderReq{}
	if ok := validators.Validate(c, &req, api.PageListingOrderVal); !ok {
		return
	}
	page, err := orderService.PageListingOrder(req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页查询失败")
	} else {
		response.SuccessData(c, page)
	}
}

func (lc *OrderController) PageBySeller(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		response.ErrorStr(c, "地址不能为空")
		return
	}
	response.SuccessData(c, orderService.GetBySeller(address))
}
