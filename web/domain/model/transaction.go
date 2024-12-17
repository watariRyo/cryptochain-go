package model

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id        uuid.UUID
	OutputMap map[string]int
	Input     *Input
}

type Input struct {
	Timestamp time.Time
	Amount    int
	Address   string
	Signature *Signature
}

type Signature struct {
	R *big.Int
	S *big.Int
	// v int
}
