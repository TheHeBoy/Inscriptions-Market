package dao

import "gohub/internal/model"

type InscriptionDao struct {
	BaseDao[model.InscriptionDO]
}

var Inscription = new(InscriptionDao)
