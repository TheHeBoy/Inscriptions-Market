package api

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/page"
	"gohub/pkg/response"
)

type InscriptionController struct {
}

var inscriptionService = service.Inscription

func (mc *InscriptionController) GetLatest(c *gin.Context) {
	pageReq := page.Req{}
	if ok := validators.Validate(c, &pageReq); !ok {
		return
	}

	pageResp, err := inscriptionService.GetLatest(pageReq)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "分页查询失败")
	} else {
		response.SuccessData(c, pageResp)
	}
}
