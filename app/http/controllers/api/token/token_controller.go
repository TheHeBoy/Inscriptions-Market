package token_controller

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/validators"
	"gohub/app/services"
	"gohub/pkg/response"
)

type TokenController struct {
}

func (lc *TokenController) PageTokens(c *gin.Context) {
	pageTokensReq := PageTokensReq{}
	if ok := validators.Validate(c, &pageTokensReq, PageTokensVal); !ok {
		return
	}

	tokenService := services.TokenService{}
	pageResp := tokenService.SelectPage(pageTokensReq.Tick, pageTokensReq.PageReq)
	if pageResp.List == nil {
		response.ErrorStr(c, "Token 分页查询出错")
	} else {
		response.SuccessData(c, pageResp)
	}
}
