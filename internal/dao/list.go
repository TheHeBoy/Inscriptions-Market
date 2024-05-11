package dao

import (
	"gohub/internal/model"
)

type ListDao struct {
	BaseDao[model.ListDO]
}

var List = new(ListDao)
