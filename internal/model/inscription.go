package model

type InscriptionDO struct {
	BaseModel

	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Block       uint64 `json:"block"`
	Idx         uint32 `json:"idx"`
	Timestamp   uint64 `json:"timestamp"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

func (*InscriptionDO) TableName() string {
	return "inscriptions"
}
