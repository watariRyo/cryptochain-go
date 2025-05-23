package repository

import (
	"context"

	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

type BlockChainInterface interface {
	AddBlock(data string, tm time.TimeProvider)
	GetBlock() []*model.Block
	IsValidChain() bool
	ReplaceChain(ctx context.Context, chain []*model.Block, tm time.TimeProvider, validTransactionDataFn func(ctx context.Context, param1 []*model.Block, param2 []*model.Block) bool)
	UnmarshalAndReplaceBlock(ctx context.Context, payload []byte, tm time.TimeProvider, fn func([]*model.Block) error, validTransactionDataFn func(ctx context.Context, param1 []*model.Block, param2 []*model.Block) bool)
}
