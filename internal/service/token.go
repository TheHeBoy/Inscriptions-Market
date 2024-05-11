package service

import (
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/internal/request"
)

type TokenService struct {
}

var Token = new(TokenService)
var tokenDao = dao.Token
var holderDao = dao.Holder
var msc20Dao = dao.Msc20

func (*TokenService) SelectPage(tick string, pageReq request.PageReq) (*request.PageResp[model.TokenDO], error) {
	pageReq.Fields = append(pageReq.Fields, "deploy_at")
	pageReq.Orders = append(pageReq.Orders, "desc")

	return tokenDao.Model().SelectPage(pageReq).
		WhereIf(tick != "", "Tick like ?", "%"+tick+"%").
		Page()
}

type GetTokensByAddressResp struct {
	Tick            string `json:"tick"`
	HoldersNum      uint64 `json:"holdersNum"`
	InscriptionsNum uint64 `json:"inscriptionsNum"`
}

func (*TokenService) GetTokensByAddress(address string) []GetTokensByAddressResp {
	var holder []model.HolderDO
	holderDao.Model().Where("address = ?", address).Find(&holder)

	ticks := make([]any, len(holder))
	for i := range holder {
		ticks[i] = holder[i].Tick
	}

	resp := make([]GetTokensByAddressResp, 0)

	var msc20s []model.Msc20DO
	msc20Dao.Model().
		Where("valid = ?", 1).
		Where("Tick in ?", ticks...).
		Where("operation = ?", "deploy").
		Find(&msc20s)

	if len(msc20s) == 0 {
		return resp
	}

	inscriptionsNumsMap := make(map[string]uint64)
	for _, inscription := range msc20s {
		inscriptionsNumsMap[inscription.Tick] = inscription.ID
	}

	for i := range holder {
		h := holder[i]
		resp = append(resp, GetTokensByAddressResp{
			Tick:            h.Tick,
			HoldersNum:      h.Amount,
			InscriptionsNum: inscriptionsNumsMap[h.Tick],
		})
	}

	return resp
}
