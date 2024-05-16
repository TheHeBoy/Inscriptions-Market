package api

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/request"
	"gohub/internal/request/validators"
	"gohub/internal/service"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

type InscriptionController struct {
}

var inscriptionService = service.Inscription

func (mc *InscriptionController) GetLatest(c *gin.Context) {
	pageReq := request.PageReq{}
	if ok := validators.Validate(c, &pageReq, validators.PageReqVal); !ok {
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
