package bigint

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"math/big"
)

type BigInt struct {
	*big.Int
}

func New() *BigInt {
	return &BigInt{big.NewInt(0)}
}

func addInt(x, y *BigInt) {
	if x.Int == nil {
		x.Int = big.NewInt(0)
	}
	if y.Int == nil {
		y.Int = big.NewInt(0)
	}
}

func (x *BigInt) Add(y *BigInt) {
	addInt(x, y)
	x.Int.Add(x.Int, y.Int)
}

func (x *BigInt) Sub(y *BigInt) {
	addInt(x, y)
	x.Int.Sub(x.Int, y.Int)
}

func (x *BigInt) Cmp(y *BigInt) int {
	addInt(x, y)
	return x.Int.Cmp(y.Int)
}

func (b *BigInt) Scan(value interface{}) error {
	b.Int = big.NewInt(0)
	switch v := value.(type) {
	case string:
	case []byte:
		b.SetString(string(v), 10)
	default:
		return fmt.Errorf("bigint: cannot convert %T to BigInt", value)
	}
	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	if b.String() == "<nil>" {
		return "0", nil
	}
	return b.String(), nil
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	if b.String() == "<nil>" {
		return json.Marshal("0")
	}
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	b.Int = big.NewInt(0)
	b.SetString(str, 10)
	return nil
}

func (b *BigInt) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "varchar(255)"
	}
	return ""
}
