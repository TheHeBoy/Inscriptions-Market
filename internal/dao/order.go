package dao

import (
	"gohub/internal/enum"
	"gohub/internal/model"
	"gorm.io/gorm"
)

type OrderDao struct {
	BaseDao[model.OrderDO]
}

var Order = new(OrderDao)

func (dao *OrderDao) Tx(db *gorm.DB) *OrderDao {
	return &OrderDao{BaseDao[model.OrderDO]{DB: db}}
}

func (dao *OrderDao) ExistByListHash(listHash string) *model.OrderDO {
	return dao.Model().Where("list_hash = ?", listHash).Exist()
}

func (dao *OrderDao) UpdateStatusByListHash(listHash string, status enum.OrderStatus) error {
	return dao.Model().Where("list_hash = ?", listHash).Update("status", status.Code).Error
}

type OrderLogDao struct {
	BaseDao[model.OrderLogDO]
}

var OrderLog = new(OrderLogDao)

func (dao *OrderLogDao) Tx(db *gorm.DB) *OrderLogDao {
	return &OrderLogDao{BaseDao[model.OrderLogDO]{DB: db}}
}

type ListDao struct {
	BaseDao[model.ListDO]
}

var List = new(ListDao)
