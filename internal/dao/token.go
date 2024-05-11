package dao

import (
	"gohub/internal/model"
)

type TokenDao struct {
	BaseDao[model.TokenDO]
}

var Token = new(TokenDao)

type HolderDao struct {
	BaseDao[model.HolderDO]
}

var Holder = new(HolderDao)

type GetByAddressResp struct {
	Collected []any `json:"collected"`
	Deploy    []any `json:"deploy"`
}
