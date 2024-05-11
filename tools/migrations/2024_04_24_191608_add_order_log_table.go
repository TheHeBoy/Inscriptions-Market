package migrations

import (
	"database/sql"
	"fmt"
	"gohub/internal/model"
	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

type OrderLogDO struct {
	model.BaseModel

	Address     string   `gorm:"column:address" json:"address"`
	Topics      []string `gorm:"serializer:json" json:"topics"`
	Data        string   `gorm:"column:data" json:"data"`
	BlockNumber int64    `gorm:"column:block_number" json:"blockNumber"`
	TxHash      string   `gorm:"column:data"  json:"txHash"`
	TxIndex     uint     `gorm:"column:tx_index" json:"txIndex"`
	Index       uint     `gorm:"column:index" json:"index"`
	Status      int      `gorm:"column:status" json:"status"`

	model.CommonTimestampsField
}

func (*OrderLogDO) TableName() string {
	return "order_log"
}

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		err := migrator.AutoMigrate(&OrderLogDO{})
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		err := migrator.DropTable(&OrderLogDO{})
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	}

	migrate.Add("2024_04_24_191608_add_order_log_table", up, down)
}
