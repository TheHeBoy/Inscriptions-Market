package service

import (
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/enum"
	"gohub/internal/model"
	"gohub/internal/request"
	"gohub/internal/request/api"
	"gohub/pkg/logger"
)

type OrderService struct {
}

var Order = new(OrderService)
var orderDao = dao.Order

func (*OrderService) Create(req api.CreateOrderReq) error {
	orderDO := orderDao.ExistByListHash(req.ListHash)
	// 已经发送了 list 铭文到合约里, 但是还没有签名
	if orderDO != nil && orderDO.Status == enum.OrderStatusWaitSignEnum.Code && orderDO.Tick == req.Tick &&
		orderDO.Seller == req.Seller && orderDO.Amount == req.Amount {
		// 添加价格和手续费
		orderDO.Price = req.Price
		orderDO.CreatorFeeRate = req.CreatorFeeRate
		orderDO.Signature = req.Signature
		orderDao.Model().Where("list_hash = ?", req.ListHash).Save(orderDO)
		return nil
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

func (*OrderService) PageListingOrder(req api.PageListingOrderReq) (*request.PageResp[model.OrderDO], error) {
	// 查询出所有状态为上架的订单
	return orderDao.Model().SelectPage(req.PageReq).
		Where("status = ?", enum.OrderStatusListingEnum.Code).
		WhereIf(req.Tick != "", "Tick like ?", "%"+req.Tick+"%").
		Page()
}

func (*OrderService) GetBySeller(address string) []model.OrderDO {
	var orderDOs []model.OrderDO
	orderDao.Model().Where("seller = ?", address).Find(&orderDOs)
	return orderDOs
}

func (*OrderService) Listing(list model.ListDO) {
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
