package service

import (
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/enum"
	"gohub/internal/model"
	"gohub/internal/request"
	"gohub/pkg/bigint"
	"gohub/pkg/config"
	"sort"
	"time"
)

type TokenService struct {
}

var Token = new(TokenService)
var tokenDao = dao.Token
var holderDao = dao.Holder

func (*TokenService) PageTokens(tick string, pageReq request.PageReq) (*request.PageResp[model.TokenDO], error) {
	pageReq.Fields = append(pageReq.Fields, "deploy_at")
	pageReq.Orders = append(pageReq.Orders, "desc")

	return tokenDao.Model().SelectPage(pageReq).
		WhereIf(tick != "", "Tick like ?", "%"+tick+"%").
		Page()
}

type PageListingTokenResp struct {
	Tick        string        `json:"tick"`
	FloorPrice  bigint.BigInt `json:"floorPrice"`
	VolumeDay   bigint.BigInt `json:"volumeDay"`
	SalesDay    int           `json:"salesDay"`
	Owners      int           `json:"owners"`
	TotalVolume bigint.BigInt `json:"totalVolume"`
	TotalSales  int           `json:"totalSales"`
	Listed      int           `json:"listed"`
}

func (*TokenService) PageListingToken(tick string, pageReq request.PageReq) (*request.PageResp[PageListingTokenResp], error) {
	var orderDOs []model.OrderDO
	orderDao.Model().WhereIf(tick != "", "tick like ?", "%"+tick+"%").Find(&orderDOs)
	var isInPastDay = func(t time.Time) bool {
		return t.After(time.Now().Add(-24 * time.Hour))
	}

	// calculate the data
	resp := make(map[string]*PageListingTokenResp)
	for i := range orderDOs {
		o := orderDOs[i]
		r := resp[o.Tick]

		if o.Status == enum.OrderStatusListingEnum.Code {
			if r != nil {
				r.Listed++
			} else {
				resp[o.Tick] = &PageListingTokenResp{
					Tick:   o.Tick,
					Listed: 1,
				}
			}
		} else if o.Status == enum.OrderStatusSoldEnum.Code {
			var volumeDay bigint.BigInt
			var salesDay int
			var price = o.Price
			if isInPastDay(o.CreatedAt) {
				volumeDay = price
				salesDay = 1
			}

			if r != nil {
				if price.Cmp(&r.FloorPrice) > 0 {
					r.FloorPrice = price
				}
				r.TotalVolume.Add(&price)
				r.TotalSales++
				r.VolumeDay.Add(&volumeDay)
				r.SalesDay += salesDay
			} else {
				resp[o.Tick] = &PageListingTokenResp{
					Tick:        o.Tick,
					FloorPrice:  price,
					VolumeDay:   volumeDay,
					SalesDay:    salesDay,
					TotalVolume: price,
					TotalSales:  1,
				}
			}
		}
	}

	rows, err := holderDao.Model().Select("tick, count(tick) as num").Where("amount > 0").Group("tick").Rows()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tick string
		var num int
		err := rows.Scan(&tick, &num)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if _, ok := resp[tick]; ok {
			resp[tick].Owners = num
		}
	}

	// sort by totalVolume
	respSlice := make([]PageListingTokenResp, 0, len(resp))
	for _, value := range resp {
		respSlice = append(respSlice, *value)
	}
	sort.Slice(respSlice, func(i, j int) bool {
		return respSlice[i].TotalVolume.Cmp(&respSlice[j].TotalVolume) < 0
	})

	// page operation
	if pageReq.PageSize <= 0 {
		pageReq.PageSize = config.GetInt("page.page_size")
	}
	if pageReq.PageNo <= 0 {
		pageReq.PageNo = 1
	}
	start := (pageReq.PageNo - 1) * pageReq.PageSize
	end := start + pageReq.PageSize
	if end > len(respSlice) {
		end = len(respSlice)
	}
	page := respSlice[start:end]

	return &request.PageResp[PageListingTokenResp]{
		Total: int64(len(resp)),
		List:  page,
	}, nil
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
		if h.Amount > 0 {
			resp = append(resp, GetTokensByAddressResp{
				Tick:            h.Tick,
				HoldersNum:      h.Amount,
				InscriptionsNum: inscriptionsNumsMap[h.Tick],
			})
		}
	}

	return resp
}
