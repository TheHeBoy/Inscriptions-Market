package migrations

import (
	"database/sql"
	"gohub/internal/model"
	"gohub/pkg/logger"
	"gohub/pkg/migrate"
	"gorm.io/gorm"
)

// SignOrder 签名的数据结构
type SignOrder struct {
	Seller         string `gorm:"column:seller" json:"seller"`
	ListHash       string `gorm:"column:list_hash;size:66;unique" json:"listHash"`
	Tick           string `gorm:"column:tick;size:18" json:"tick"`
	Amount         int    `gorm:"column:amount" json:"amount"`
	Price          int    `gorm:"column:price" json:"price"`
	CreatorFeeRate int    `gorm:"column:creator_fee_rate" json:"creatorFeeRate"`
}

type OrderDO struct {
	model.BaseModel

	SignOrder
	Signature string `gorm:"column:signature" json:"signature"`
	Status    string `gorm:"column:status;type:enum('Listing', 'Sold','Unfounded' ,'Cancelled')" json:"status"`

	model.CommonTimestampsField
}

func (*OrderDO) TableName() string {
	return "orders"
}

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		err := migrator.AutoMigrate(&OrderDO{})
		if err != nil {
			logger.Error(err)
		}
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		err := migrator.DropTable(&OrderDO{})
		if err != nil {
			logger.Error(err)
		}
	}

	migrate.Add("2024_04_22_010219_add_order_table", up, down)
}
