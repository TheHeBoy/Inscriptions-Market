package service

import (
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/enum"
	"gohub/internal/model"
	"gohub/pkg/bigint"
	"gohub/pkg/config"
	"gohub/pkg/page"
	"sort"
	"time"
)

type TokenService struct {
}

var Token = new(TokenService)
var tokenDao = dao.Token
var holderDao = dao.Holder

func (*TokenService) PageTokens(tick string, pageReq page.Req) (*page.Resp[model.TokenDO], error) {
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

func (*TokenService) PageListingToken(tick string, pageReq page.Req) (*page.Resp[PageListingTokenResp], error) {
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
	list := respSlice[start:end]

	return &page.Resp[PageListingTokenResp]{
		Total: int64(len(resp)),
		List:  list,
	}, nil
}

type GetTokensByAddressResp struct {
	Tick            string `json:"tick"`
	HoldersNum      uint64 `json:"holdersNum"`
	InscriptionsNum uint64 `json:"inscriptionsNum"`
}

func (*TokenService) GetTokensByAddress(address string) ([]GetTokensByAddressResp, error) {
	var holder []model.HolderDO
	holderDao.Model().Where("address = ?", address).Find(&holder)

	ticks := make([]string, len(holder))
	for i := range holder {
		ticks[i] = holder[i].Tick
	}

	resp := make([]GetTokensByAddressResp, 0)

	rows, err := msc20Dao.Model().
		Select("tick", "id").
		Where("valid = ?", 1).
		Where("tick in ?", ticks).
		Where("operation = ?", "deploy").
		Rows()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	inscriptionsNumsMap, err := dao.MapRows[string, uint64](rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(inscriptionsNumsMap) == 0 {
		return resp, nil
	}

	// indexer 还没有收到 list 时，减少持有者数量
	rows, err = orderDao.Model().
		Select("tick", "amount").
		Where("tick in ?", ticks).
		Where("seller = ?", address).
		Where("status = ?", enum.OrderStatusWaitListEnum.Code).Rows()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderAmountMap, err := dao.MapRows[string, uint64](rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range holder {
		h := holder[i]
		if h.Amount > 0 {
			resp = append(resp, GetTokensByAddressResp{
				Tick:            h.Tick,
				HoldersNum:      h.Amount - orderAmountMap[h.Tick],
				InscriptionsNum: inscriptionsNumsMap[h.Tick],
			})
		}
	}

	return resp, nil
}
