// Package requests 处理请求数据和表单验证
package requests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/errorcode"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

// ValidatorFunc 验证函数类型
type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

// Validate 控制器里调用示例：
//
//	if ok := requests.Validate(c, &requests.UserSaveRequest{}, requests.UserSave); !ok {
//	    return
//	}
func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {

	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(obj); err != nil {
		logger.Error(err)
		response.ErrorCustom(c, errorcode.BAD_REQUEST.Code, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}

	// 2. 表单验证
	errs := handler(obj, c)

	jsonData, err := json.Marshal(errs)
	if err != nil {
		logger.Error(err)
		response.Error(c, errorcode.JSON_FORMAT_ERROR)
		return false
	}

	// 3. 判断验证是否通过
	if len(errs) > 0 {
		response.ErrorCustom(c, errorcode.BAD_REQUEST.Code, fmt.Sprintf("请求参数不正确:%s", jsonData))
		return false
	}

	return true
}

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
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
