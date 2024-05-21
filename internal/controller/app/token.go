package app

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/app"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type TokenController struct {
}

var tokenService = service.Token

func (lc *TokenController) PageTokens(c *gin.Context) {
	pageTokensReq := &app.PageTokensReq{}
	if ok := validators.Validate(c, pageTokensReq); !ok {
		return
	}

	pageResp, err := tokenService.PageTokens(pageTokensReq.Tick, pageTokensReq.Req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页查询失败")
	} else {
		response.SuccessData(c, pageResp)
	}
}

func (lc *TokenController) PageListingToken(c *gin.Context) {
	pageTokensReq := app.PageTokensReq{}
	if ok := validators.Validate(c, &pageTokensReq); !ok {
		return
	}

	pageResp, err := tokenService.PageListingToken(pageTokensReq.Tick, pageTokensReq.Req)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页查询失败")
	} else {
		response.SuccessData(c, pageResp)
	}
}

func (lc *TokenController) GetTokensByAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		response.ErrorStr(c, "地址不能为空")
		return
	}
	resp, err := tokenService.GetTokensByAddress(address)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "查询Token失败")
		return
	} else {
		response.SuccessData(c, resp)
	}
}
