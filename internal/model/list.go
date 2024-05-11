package model

type ListDO struct {
	BaseModel

	Hash     string `gorm:"column:hash" json:"hash"`
	Owner    string `gorm:"column:owner" json:"owner"`
	Exchange string `gorm:"column:exchange" json:"exchange"`
	Tick     string `gorm:"column:tick" json:"tick"`
	Amount   uint64 `gorm:"column:amount" json:"amount"`
}

func (*ListDO) TableName() string {
	return "lists"
}
