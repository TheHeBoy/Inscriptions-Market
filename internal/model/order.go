package model

type OrderDO struct {
	BaseModel

	SignOrder
	Signature string `gorm:"column:signature" json:"signature"`
	Status    string `gorm:"column:status" json:"status"`

	CommonTimestampsField
}

// SignOrder 签名的数据结构
type SignOrder struct {
	Seller         string `gorm:"column:seller" json:"seller"`
	ListHash       string `gorm:"column:list_hash;unique" json:"listHash"`
	Tick           string `gorm:"column:tick;" json:"tick"`
	Amount         uint64 `gorm:"column:amount" json:"amount"`
	Price          int    `gorm:"column:price" json:"price"`
	CreatorFeeRate int    `gorm:"column:creator_fee_rate" json:"creatorFeeRate"`
}

func (*OrderDO) TableName() string {
	return "orders"
}

type OrderLogDO struct {
	BaseModel

	Address     string   `gorm:"column:address" json:"address"`
	Topics      []string `gorm:"serializer:json" json:"topics"`
	Data        string   `gorm:"column:data" json:"data"`
	BlockNumber int64    `gorm:"column:block_number" json:"blockNumber"`
	TxHash      string   `gorm:"column:data"  json:"txHash"`
	TxIndex     uint     `gorm:"column:tx_index" json:"txIndex"`
	Index       uint     `gorm:"column:index" json:"index"`
	Status      string   `gorm:"column:status" json:"status"`

	CommonTimestampsField
}

func (*OrderLogDO) TableName() string {
	return "order_log"
}
