package app

import (
	"github.com/thedevsaddam/govalidator"
	"gohub/internal/request/validators"
	"gohub/pkg/bigint"
)

type CreateOrderReq struct {
	Seller         string        `json:"seller"`
	ListHash       string        `json:"listHash"`
	Tick           string        `json:"tick"`
	Amount         uint64        `json:"amount"`
	Price          bigint.BigInt `json:"price"`
	CreatorFeeRate int           `json:"creatorFeeRate"`
	Signature      string        `json:"signature"`
}

func (r *CreateOrderReq) Validator() map[string][]string {
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

	errs := validators.ValidateData(r, rules, messages)

	return errs
}

type SignOrderReq struct {
	ID             uint64        `json:"id"`
	Price          bigint.BigInt `json:"price"`
	CreatorFeeRate int           `json:"creatorFeeRate"`
	Signature      string        `json:"signature"`
}

func (r *SignOrderReq) Validator() map[string][]string {
	return make(map[string][]string)
}

type GetListingOrderByTickReq struct {
	Tick string `json:"tick" form:"tick"`
}

func (r *GetListingOrderByTickReq) Validator() map[string][]string {
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

	return validators.ValidateData(r, rules, messages)
}
