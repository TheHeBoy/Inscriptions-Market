package services

import (
	"gohub/app/models/token"
	"gohub/pkg/dal/database"
)

type TokenService struct {
}

func (*TokenService) SelectPage(tick string, pageReq database.PageReq) *database.PageResp[token.Token] {
	pageReq.Fields = append(pageReq.Fields, "deploy_at")
	pageReq.Orders = append(pageReq.Orders, "desc")

	db := database.NewPagePaginate(pageReq)
	db.Model(&token.Token{})
	if tick != "" {
		db = db.Order("deploy_at desc").Where("tick like ?", "%"+tick+"%")
	}

	var pageResp = new(database.PageResp[token.Token])
	db.Count(&pageResp.Total)
	db.Find(&pageResp.List)
	return pageResp
}
