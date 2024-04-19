// Package validators 存放自定义规则和验证器
package validators

import (
	"fmt"
	"gohub/pkg/config"
	"gohub/pkg/dal/database"
)

// ValidatePage 自定义规则，验证分页参数
func ValidatePage(pageReq database.PageReq, errs map[string][]string) map[string][]string {

	maxPageSize := config.GetInt("page.max_page_size")
	if pageReq.PageSize > maxPageSize {
		errs["pageNo"] = append(errs["pageNo"], fmt.Sprintf("页码超出最大限制%d", maxPageSize))
	}

	fields := pageReq.Fields
	orders := pageReq.Orders

	if len(fields) != len(orders) {
		errs["fields"] = append(errs["fields"], "fields 和 orders 长度不一致")
	}

	for _, field := range fields {
		if field == "" {
			errs["fields"] = append(errs["fields"], "fields 不能为空")
			break
		}
	}

	for _, order := range orders {
		if order != "asc" && order != "desc" {
			errs["orders"] = append(errs["orders"], "order 只能是 asc 或 desc")
			break
		}
	}
	return errs
}
