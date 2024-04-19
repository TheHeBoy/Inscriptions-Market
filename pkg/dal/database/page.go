package database

import (
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gorm.io/gorm"
)

type PageReq struct {
	PageNo   int      `json:"pageNo" valid:"pageNo" form:"pageNo"`
	PageSize int      `json:"pageSize" valid:"pageSize" form:"pageSize"`
	Orders   []string `json:"orders" valid:"orders" form:"orders"`
	Fields   []string `json:"fields" valid:"fields" form:"fields"`
}

type PageResp[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"` // nil 表示查询出错， [] 表示数据为空
}

type Conditions func(*gorm.DB) *gorm.DB

func paginate(pageReq PageReq) Conditions {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := pageReq.PageNo
		if pageNo <= 0 {
			pageNo = 1
		}

		pageSize := pageReq.PageSize
		maxPageSize := config.GetInt("page.max_page_size")
		switch {
		case pageSize > maxPageSize:
			logger.Warn("pageSize is too large, set to maxPageSize")
			pageSize = maxPageSize
		case pageSize <= 0:
			logger.Warn("pageSize is invalid, set to default")
			pageSize = config.GetInt("page.page_size")
		}

		offset := (pageNo - 1) * pageSize

		// Handle Orders
		fields := pageReq.Fields
		orders := pageReq.Orders
		if len(fields) != len(orders) {
			logger.Warn("fields and orders are not matched")
		}
		for i := 0; i < len(fields); i++ {
			db = db.Order(fields[i] + " " + orders[i])

		}
		return db.Offset(offset).Limit(pageSize)
	}
}

func NewPagePaginate(pageReq PageReq) *gorm.DB {
	return DB.Scopes(paginate(pageReq))
}
