// Package daos 对于Gorm的封装增强，以及分页查询的封装。
package dao

import (
	"github.com/pkg/errors"
	"gohub/internal/request"
	"gohub/pkg/config"
	"gohub/pkg/database"
	"gohub/pkg/logger"
	"gorm.io/gorm"
)

type BaseDao[T any] struct {
	*gorm.DB
}

func (dao *BaseDao[T]) New() *BaseDao[T] {
	return &BaseDao[T]{DB: database.DB}
}

func (dao *BaseDao[T]) Model() *BaseDao[T] {
	var model = new(T)
	baseDao := dao.New()
	baseDao.DB = database.DB.Model(model)
	return baseDao
}

func paginate(pageReq request.PageReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := pageReq.PageNo
		if pageNo <= 0 {
			pageNo = 1
		}

		pageSize := pageReq.PageSize
		maxPageSize := config.GetInt("page.max_page_size")
		switch {
		case pageSize > maxPageSize:
			pageSize = maxPageSize
		case pageSize <= 0:
			pageSize = config.GetInt("page.page_size")
		}

		offset := (pageNo - 1) * pageSize

		// Handle Orders
		fields := pageReq.Fields
		orders := pageReq.Orders

		for i := 0; i < len(fields); i++ {
			db = db.Order(fields[i] + " " + orders[i])

		}
		return db.Offset(offset).Limit(pageSize)
	}
}

func (dao *BaseDao[T]) SelectPage(pageReq request.PageReq) *BaseDao[T] {
	*dao.DB = *dao.DB.Scopes(paginate(pageReq))
	return dao
}

func (dao *BaseDao[T]) Where(query interface{}, args ...interface{}) *BaseDao[T] {
	*dao.DB = *dao.DB.Where(query, args)
	return dao
}

func (dao *BaseDao[T]) WhereIf(condition bool, query interface{}, args ...interface{}) *BaseDao[T] {
	if condition {
		*dao.DB = *dao.DB.Where(query, args)
		return dao
	}
	return dao
}

func (dao *BaseDao[T]) Page() (*request.PageResp[T], error) {
	var pageResp = new(request.PageResp[T])
	err := dao.DB.Count(&pageResp.Total).Error
	if err != nil {
		return nil, err
	}

	err = dao.DB.Find(&pageResp.List).Error
	if err != nil {
		return nil, err
	}
	return pageResp, nil
}

// Exist
//
//	@Description: 判断是否存在
//	@receiver dao
//	@return *T 如果存在返回实体，否则返回nil
func (dao *BaseDao[T]) Exist() *T {
	var model T
	err := dao.DB.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		logger.Error(err)
		return nil
	} else {
		return &model
	}
}
