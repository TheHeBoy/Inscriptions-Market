package token

import "gohub/app/models"

type Token struct {
	models.BaseModel

	Tick        string `gorm:"column:tick;size:20;unique;comment:代币名称" json:"tick"`
	Max         uint64 `gorm:"column:max;comment:总mint最大值" json:"max"`
	Limit       uint64 `gorm:"column:limit;comment:每次mint的最大值" json:"limit"`
	Minted      uint64 `gorm:"column:minted;comment:已经mint的数量" json:"minted"`
	Progress    string `gorm:"column:progress;size:5;comment:mint进度" json:"progress"`
	Trxs        uint32 `gorm:"column:trxs;comment:操作交易数" json:"trxs"`
	CompletedAt uint64 `gorm:"column:completed_at;comment:mint全部完成的时间" json:"completedAt"`
	DeployAt    uint64 `gorm:"column:deploy_at;index" json:"deployAt"`
}
