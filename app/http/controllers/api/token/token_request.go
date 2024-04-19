package token_controller

import (
	"github.com/thedevsaddam/govalidator"
	"gohub/app/http/validators"
	"gohub/pkg/dal/database"
)

type PageTokensReq struct {
	database.PageReq
	Tick string `json:"tick" valid:"tick" form:"tick"`
}

func PageTokensVal(data any) map[string][]string {
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
	_data := data.(*PageTokensReq).PageReq
	errs = validators.ValidatePage(_data, errs)

	return errs
}
