package api

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/api"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type TokenController struct {
}

var tokenService = service.Token

func (lc *TokenController) PageTokens(c *gin.Context) {
	pageTokensReq := api.PageTokensReq{}
	if ok := validators.Validate(c, &pageTokensReq, api.PageTokensVal); !ok {
		return
	}

	pageResp, err := tokenService.PageTokens(pageTokensReq.Tick, pageTokensReq.PageReq)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页查询失败")
	} else {
		response.SuccessData(c, pageResp)
	}
}

func (lc *TokenController) PageListingToken(c *gin.Context) {
	pageTokensReq := api.PageTokensReq{}
	if ok := validators.Validate(c, &pageTokensReq, api.PageTokensVal); !ok {
		return
	}

	pageResp, err := tokenService.PageListingToken(pageTokensReq.Tick, pageTokensReq.PageReq)
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

	response.SuccessData(c, tokenService.GetTokensByAddress(address))
}
