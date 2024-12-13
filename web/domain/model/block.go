package model

type Block struct {
	Timestamp  string
	LastHash   string
	Hash       string
	Difficulty int
	Nonce      int
	Data       string
}