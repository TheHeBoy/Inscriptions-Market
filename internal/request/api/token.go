package api

import (
	"github.com/thedevsaddam/govalidator"
	"gohub/internal/request/validators"
	"gohub/pkg/page"
)

type PageTokensReq struct {
	page.Req
	Tick string `json:"tick" valid:"tick" form:"tick"`
}

func (r *PageTokensReq) Validator() map[string][]string {
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

	// 添加分页校验
	_data := &r.Req
	errs = page.ValidatePage(_data, errs)

	return errs
}
