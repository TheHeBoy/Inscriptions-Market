package service

import (
	"gohub/internal/dao"
	"gohub/internal/model"
)

type Msc20Service struct {
}

var Msc20 = new(Msc20Service)
var msc20Dao = dao.Msc20

type GetMsc20ByAddressResp struct {
	ID        uint64 `json:"id"`
	Content   string `json:"content"`
	Hash      string `json:"hash"`
	Operation string `json:"operation"`
	Timestamp uint64 `json:"timestamp"`
}

func (s *Msc20Service) GetMsc20ByAddress(address string) []GetMsc20ByAddressResp {
	var msc20s []model.Msc20DO
	msc20Dao.Model().Where("valid = ?", 1).Where("`from` = ?", address).Find(&msc20s)

	txs := make([]any, len(msc20s))
	for i := range msc20s {
		txs[i] = msc20s[i].Hash
	}

	var insDos []model.InscriptionDO
	inscriptionDao.Model().Where("hash in ?", txs...).Find(&insDos)

	insDoMap := make(map[string]string)
	for i := range insDos {
		insDoMap[insDos[i].Hash] = insDos[i].Content
	}

	resp := make([]GetMsc20ByAddressResp, len(msc20s))
	for i := range msc20s {
		resp[i] = GetMsc20ByAddressResp{
			ID:        msc20s[i].ID,
			Operation: msc20s[i].Operation,
			Hash:      msc20s[i].Hash,
			Timestamp: msc20s[i].Timestamp,
			Content:   insDoMap[msc20s[i].Hash],
		}
	}
	return resp
}
