package api

import (
	"github.com/thedevsaddam/govalidator"
	"gohub/internal/request/validators"
)

type GetMessageReq struct {
	Address string `json:"address" form:"address"` // 账户地址
}

func (r *GetMessageReq) Validator() map[string][]string {
	rules := govalidator.MapData{
		"address": []string{"required"},
	}

	messages := govalidator.MapData{
		"address": []string{
			"required:账户地址为必填项",
		},
	}
	return validators.ValidateData(r, rules, messages)
}

type LoginBySignatureReq struct {
	Address   string `json:"address"`   // 账户地址
	Signature string `json:"signature"` // 签名
}

func (r *LoginBySignatureReq) Validator() map[string][]string {

	rules := govalidator.MapData{
		"address":   []string{"required"},
		"signature": []string{"required"},
	}

	messages := govalidator.MapData{
		"address": []string{
			"required:账户地址为必填项",
		},
		"signature": []string{
			"required:签名为必填项",
		},
	}

	return validators.ValidateData(r, rules, messages)
}
