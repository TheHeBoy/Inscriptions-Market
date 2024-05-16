package service

import (
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/enum"
	"gohub/internal/model"
	"gohub/internal/request/api"
	"gohub/pkg/bigint"
	"gohub/pkg/logger"
)

type OrderService struct {
}

var Order = new(OrderService)
var orderDao = dao.Order
var listDao = dao.List

func (s *OrderService) Create(req api.CreateOrderReq) error {
	orderDO := orderDao.ExistByListHash(req.ListHash)
	// 已经发送了 list 铭文到合约里, 但是还没有签名
	if orderDO != nil {
		if orderDO.Tick == req.Tick &&
			orderDO.Seller == req.Seller && orderDO.Amount == req.Amount {
			return errors.New("order info is error")
		}
		return s.sign(orderDO, req.Price, req.CreatorFeeRate, req.Signature)
	}

	// 创建订单
	order := model.OrderDO{
		SignOrder: model.SignOrder{
			Seller:         req.Seller,
			ListHash:       req.ListHash,
			Tick:           req.Tick,
			Amount:         req.Amount,
			Price:          req.Price,
			CreatorFeeRate: req.CreatorFeeRate,
		},
		Status:    enum.OrderStatusWaitListEnum.Code,
		Signature: req.Signature,
	}

	if err := orderDao.Model().Create(&order).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *OrderService) sign(orderDO *model.OrderDO, price bigint.BigInt, creatorFeeRate int, signature string) error {
	if orderDO.Status != enum.OrderStatusWaitSignEnum.Code {
		return errors.New("order status is error, need wait sign status")
	}

	// 添加价格和手续费
	orderDO.Price = price
	orderDO.CreatorFeeRate = creatorFeeRate
	orderDO.Signature = signature
	orderDO.Status = enum.OrderStatusListingEnum.Code
	if err := orderDao.Model().Where("id = ?", orderDO.ID).Save(orderDO).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *OrderService) SignOrder(req api.SignOrderReq) error {
	orderDO := orderDao.Model().Where("id = ?", req.ID).Exist()
	if orderDO == nil {
		return errors.New("order not exist")
	}
	return s.sign(orderDO, req.Price, req.CreatorFeeRate, req.Signature)
}

func (*OrderService) GetListingOrderByTick(tick string) []model.OrderDO {
	var orderDOs []model.OrderDO
	// 查询出所有状态为上架的订单
	orderDao.Model().
		Where("status = ?", enum.OrderStatusListingEnum.Code).
		Where("Tick = ?", tick).
		Find(&orderDOs)
	return orderDOs
}

func (*OrderService) GetBySeller(address string) []model.OrderDO {
	var orderDOs []model.OrderDO
	orderDao.Model().Where("seller = ?", address).Find(&orderDOs)
	return orderDOs
}

func (*OrderService) List(list model.ListDO) {
	orderDO := orderDao.ExistByListHash(list.Hash)
	if orderDO == nil {
		// 创建订单
		order := model.OrderDO{
			SignOrder: model.SignOrder{
				Seller:   list.Owner,
				ListHash: list.Hash,
				Tick:     list.Tick,
				Amount:   list.Amount,
			},
			Status: enum.OrderStatusWaitSignEnum.Code,
		}

		if err := orderDao.Model().Create(&order).Error; err != nil {
			logger.Errorv(errors.WithStack(err))
		}
	} else {
		if err := orderDao.UpdateStatusByListHash(list.Hash, enum.OrderStatusListingEnum); err != nil {
			logger.Errorv(errors.WithStack(err))
		}
	}
}

func (*OrderService) Execute(listHash string) enum.OrderLogStatus {
	return updateOrderStatus(listHash, enum.OrderStatusSoldEnum)
}

func (*OrderService) Cancel(listHash string) enum.OrderLogStatus {
	return updateOrderStatus(listHash, enum.OrderStatusCanceledEnum)
}

func updateOrderStatus(listHash string, status enum.OrderStatus) enum.OrderLogStatus {
	orderDO := orderDao.ExistByListHash(listHash)
	if orderDO == nil {
		return enum.OrderLogStatusOrderNotExist
	}
	if orderDO.Status != enum.OrderStatusListingEnum.Code {
		return enum.OrderLogStatusStatusError
	}

	err := orderDao.UpdateStatusByListHash(listHash, status)
	if err != nil {
		return enum.OrderLogStatusUpdateFailed
	}
	return enum.OrderLogStatusSuccess
}
