package service

import (
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/pkg/page"
)

type InscriptionService struct {
}

var Inscription = new(InscriptionService)
var inscriptionDao = dao.Inscription

func (s *InscriptionService) GetLatest(req page.Req) (*page.Resp[model.InscriptionDO], error) {
	return inscriptionDao.Model().Order("id desc").SelectPage(req).Page()
}
