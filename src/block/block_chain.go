package block

import (
	"context"
	"reflect"
	"time"

	"github.com/watariRyo/cryptochain-go/src/logger"
)

type BlockChain struct {
	ctx   context.Context
	Block []*Block
}

func NewBlockChain(ctx context.Context) *BlockChain {
	genesis := newGenesisBlock()
	blockChain := &BlockChain{
		ctx:   ctx,
		Block: []*Block{genesis},
	}

	return blockChain
}

func (bc *BlockChain) AddBlock(data string, timestamp time.Time) {
	lastBlock := bc.Block[len(bc.Block)-1]
	addBlock := MineBlock(lastBlock, data, timestamp)

	bc.Block = append(bc.Block, addBlock)
}

func (bc *BlockChain) IsValidChain() bool {
	genesis := newGenesisBlock()
	if !reflect.DeepEqual(bc.Block[0], genesis) {
		return false
	}

	actualLastHash := genesis.Hash
	for _, block := range bc.Block[1:] {
		if actualLastHash != block.LastHash {
			return false
		}
		validatedHash := cryptoHash(block.Timestamp.String(), block.LastHash, block.Data)
		if block.Hash != validatedHash {
			return false
		}
		actualLastHash = block.Hash
	}

	return true
}

func (bc *BlockChain) ReplaceChain(chain *BlockChain) {
	if len(chain.Block) <= len(bc.Block) {
		logger.Errorf(bc.ctx, "The incoming chain must be longer.")
		return
	}
	if !chain.IsValidChain() {
		logger.Errorf(bc.ctx, "The incoming chain must be valid.")
		return
	}

	bc.Block = chain.Block
}
