package service

import "gohub/internal/dao"

type InscriptionService struct {
}

var Inscription = new(InscriptionService)
var inscriptionDao = dao.Inscription
