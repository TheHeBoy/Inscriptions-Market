package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/app/requests"
)

type GetMessageReq struct {
	Address string `json:"address" valid:"address" form:"address"` // 账户地址
}

func GetMessageVal(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"address": []string{"required"},
	}

	messages := govalidator.MapData{
		"address": []string{
			"required:账户地址为必填项",
		},
	}
	return requests.ValidateData(data, rules, messages)
}

type LoginBySignatureReq struct {
	Address   string `json:"address" valid:"address"`     // 账户地址
	Signature string `json:"signature" valid:"signature"` // 签名
}

func LoginBySignatureVal(data interface{}, c *gin.Context) map[string][]string {

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

	return requests.ValidateData(data, rules, messages)
}
