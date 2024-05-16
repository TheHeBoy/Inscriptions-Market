package service

import (
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/internal/request"
)

type InscriptionService struct {
}

var Inscription = new(InscriptionService)
var inscriptionDao = dao.Inscription

func (s *InscriptionService) GetLatest(req request.PageReq) (*request.PageResp[model.InscriptionDO], error) {
	return inscriptionDao.Model().Order("id desc").SelectPage(req).Page()
}
