package model

type Transact struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}
