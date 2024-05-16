package model

type TokenDO struct {
	BaseModel

	Tick        string `gorm:"column:tick;" json:"tick"`
	Max         uint64 `gorm:"column:max" json:"max"`
	Limit       uint64 `gorm:"column:limit" json:"limit"`
	Minted      uint64 `gorm:"column:minted" json:"minted"`
	Progress    string `gorm:"column:progress" json:"progress"`
	Txs         uint32 `gorm:"column:txs" json:"txs"`
	CompletedAt uint64 `gorm:"column:completed_at" json:"completedAt"`
	DeployAt    uint64 `gorm:"column:deploy_at;index" json:"deployAt"`
}

func (*TokenDO) TableName() string {
	return "tokens"
}

type HolderDO struct {
	BaseModel

	Tick    string `gorm:"column:tick;size:20" json:"tick"`
	Address string `gorm:"column:address;size:42" json:"number"`
	Amount  uint64 `gorm:"column:amount" json:"amount"`
}

func (*HolderDO) TableName() string {
	return "holders"
}
