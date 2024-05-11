package seeders

import (
	"encoding/json"
	"gohub/pkg/console"
	"gohub/pkg/eth"
	"gohub/pkg/seed"
	"gohub/tools/factories"
	"gorm.io/gorm"
)

func init() {

	// 添加 Seeder
	seed.Add("InscriptionSeeder", func(db *gorm.DB) {
		inscriptions := factories.MakeInscription()

		for _, ins := range inscriptions {
			marshal, err := json.Marshal(ins.Ins)
			if err != nil {
				console.Error(err.Error())
			}
			trxHash, err := eth.Inscribe(ins.FormPrivateKey, ins.ToAddress, string(marshal))
			console.Success(trxHash)

			if err != nil {
				console.Error(err.Error())
			}
		}

	})
}
