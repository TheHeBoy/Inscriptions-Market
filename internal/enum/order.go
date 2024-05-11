package enum

type OrderStatus BaseEnum

var (
	OrderStatusWaitListEnum = OrderStatus{"1", "WaitList"}
	OrderStatusListingEnum  = OrderStatus{"2", "Listing"}
	OrderStatusSoldEnum     = OrderStatus{"3", "Sold"}
	OrderStatusWaitSignEnum = OrderStatus{"4", "WaitSign"} // list给了合约，但是没有签名order
	OrderStatusCanceledEnum = OrderStatus{"5", "Canceled"}
)

type OrderLogStatus BaseEnum

var (
	OrderLogStatusSuccess       = OrderLogStatus{"1", "success"}
	OrderLogStatusOrderNotExist = OrderLogStatus{"2", "Order not exist"}
	OrderLogStatusStatusError   = OrderLogStatus{"3", "Order status is error"}
	OrderLogStatusUpdateFailed  = OrderLogStatus{"4", "Order update failed"}
	OrderLogStatusDecodeFail    = OrderLogStatus{"4", "Order decode failed"}
)
