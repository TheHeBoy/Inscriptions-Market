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

	//user1 := Account{
	//	Address:    "0x90b683d65C0643fA830990787836d7E77001700E",
	//	PrivateKey: "0x41b3a95f32c9988c4f145a4758ff4647fc5fe7796f54e998e68f175ffcb677f4",
	//}
	//
	//user2 := Account{
	//	Address:    "0x96A3C77706103793Cd71acaa607817542ca08698",
	//	PrivateKey: "0x2fb15ca67c9f8eda937437324ea816a24a4c3f1fc0e5f4a0b6f4b68bba3dc98a",
	//}

	//// deploy
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

	// mint
	//insJson = append(insJson, InscriptionSeeder{
	//	FormPrivateKey: user1.PrivateKey,
	//	ToAddress:      user1.Address,
	//	Ins: map[string]string{
	//		"p":    "msc-20",
	//		"op":   "mint",
	//		"tick": "demo",
	//		"amt":  "10",
	//	},
	//})

	// transfer
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

	//// list
	//insJson = append(insJson, InscriptionSeeder{
	//	FormPrivateKey: user2.PrivateKey,
	//	ToAddress:      "0xF115fFA206aF8FC74e6D916d9945B2EC8A1f8EcC",
	//	Ins: map[string]string{
	//		"p":    "msc-20",
	//		"op":   "mint",
	//		"tick": "demo",
	//		"amt":  "10",
	//	},
	//})
	return
}
