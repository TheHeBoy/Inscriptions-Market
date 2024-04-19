// Package requests 处理请求数据和表单验证
package validators

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/response"
)

// ValidatorFunc 验证函数类型
type ValidatorFunc func(interface{}) map[string][]string

// Validate 控制器里调用示例：
//
//	if ok := requests.Validate(c, &requests.UserSaveRequest{}, requests.UserSave); !ok {
//	    return
//	}
func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {

	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(obj); err != nil {
		response.Error405(c,
			errors.Wrap(err, "请求解析错误，请确认请求格式是否正确。参数请使用 JSON 格式。"))
		return false
	}

	// 2. 表单验证
	errs := handler(obj)

	jsonData, err := json.Marshal(errs)
	if err != nil {
		response.Error(c, err)
		return false
	}

	// 3. 判断验证是否通过
	if len(errs) > 0 {
		response.Error405(c, errors.New(fmt.Sprintf("参数校验错误: %s", string(jsonData))))
		return false
	}

	return true
}

func ValidateData(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}
	// 开始验证
	return govalidator.New(opts).ValidateStruct()
}

func validateFile(c *gin.Context, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:       c.Request,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	// 调用 govalidator 的 Validate 方法来验证文件
	return govalidator.New(opts).Validate()
}
