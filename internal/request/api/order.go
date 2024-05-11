package api

import (
	"github.com/thedevsaddam/govalidator"
	"gohub/internal/request"
	"gohub/internal/request/validators"
)

type CreateOrderReq struct {
	Seller         string `json:"seller"`
	ListHash       string `json:"listHash"`
	Tick           string `json:"tick"`
	Amount         uint64 `json:"amount"`
	Price          int    `json:"price"`
	CreatorFeeRate int    `json:"creatorFeeRate"`
	Signature      string `json:"signature"`
}

func CreateOrderVal(data any) map[string][]string {
	rules := govalidator.MapData{
		"tick": []string{"min:1", "max:16"},
	}

	messages := govalidator.MapData{
		"tick": []string{
			"required:tick 为必填项",
			"min:tick 长度需大于 1",
			"max:tick 长度需小于 16",
		},
	}

	errs := validators.ValidateData(data, rules, messages)

	return errs
}

type PageListingOrderReq struct {
	request.PageReq
	Tick string `json:"tick" valid:"tick" form:"tick"`
}

func PageListingOrderVal(data any) map[string][]string {
	rules := govalidator.MapData{
		"tick": []string{"min:1", "max:16"},
	}

	messages := govalidator.MapData{
		"tick": []string{
			"required:tick 为必填项",
			"min:tick 长度需大于 1",
			"max:tick 长度需小于 16",
		},
	}

	errs := validators.ValidateData(data, rules, messages)

	// 添加分页校验
	_data := data.(*PageListingOrderReq).PageReq
	errs = validators.ValidatePage(_data, errs)
	return errs
}
