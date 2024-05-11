package dao

import (
	"gohub/internal/model"
)

type Msc20Dao struct {
	BaseDao[model.Msc20DO]
}

var Msc20 = new(Msc20Dao)
