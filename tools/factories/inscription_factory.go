// Package factories 存放工厂方法
package factories

type InscriptionSeeder struct {
	FormPrivateKey string
	ToAddress      string
	Ins            map[string]string
}

type Account struct {
	Address    string
	PrivateKey string
}

func MakeInscription() (insJson []InscriptionSeeder) {

	user1 := Account{
		Address:    "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	}

	//user2 := Account{
	//	Address:    "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
	//	PrivateKey: "0x2fb15ca67c9f8eda937437324ea816a24a4c3f1fc0e5f4a0b6f4b68bba3dc98a",
	//}

	//insJson = append(insJson, InscriptionSeeder{
	//	FormPrivateKey: user1.PrivateKey,
	//	ToAddress:      user1.Address,
	//	Ins: map[string]string{
	//		"p":    "msc-20",
	//		"op":   "deploy",
	//		"tick": "demo",
	//		"max":  "10000",
	//		"lim":  "10",
	//	},
	//})

	insJson = append(insJson, InscriptionSeeder{
		FormPrivateKey: user1.PrivateKey,
		ToAddress:      user1.Address,
		Ins: map[string]string{
			"p":    "msc-20",
			"op":   "mint",
			"tick": "demo",
			"amt":  "10",
		},
	})

	//insJson = append(insJson, InscriptionSeeder{
	//	FormPrivateKey: user1.PrivateKey,
	//	ToAddress:      user2.Address,
	//	Ins: map[string]string{
	//		"p":    "msc-20",
	//		"op":   "transfer",
	//		"tick": "demo",
	//		"amt":  "10",
	//	},
	//})

	insJson = append(insJson, InscriptionSeeder{
		FormPrivateKey: user1.PrivateKey,
		ToAddress:      "0xe7f1725e7734ce288f8367e1bb143e90bb3f0512",
		Ins: map[string]string{
			"p":    "msc-20",
			"op":   "list",
			"tick": "demo",
			"amt":  "2",
		},
	})
	return
}
