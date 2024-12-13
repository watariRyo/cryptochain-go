package repository

import (
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

type BlockChainInterface interface {
	AddBlock(data string, tm time.TimeProvider)
	GetBlock() []*model.Block
	IsValidChain() bool
	ReplaceChain(chain []*model.Block, tm time.TimeProvider)
	UnmarshalAndReplaceBlock(payload []byte, tm time.TimeProvider)
}