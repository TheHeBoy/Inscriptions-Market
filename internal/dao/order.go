package dao

import (
	"gohub/internal/enum"
	"gohub/internal/model"
)

type OrderDao struct {
	BaseDao[model.OrderDO]
}

var Order = new(OrderDao)

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

type ListDao struct {
	BaseDao[model.ListDO]
}

var List = new(ListDao)
