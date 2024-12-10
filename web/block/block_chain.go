package block

import (
	"context"
	"math"
	"reflect"
	"strconv"

	"github.com/watariRyo/cryptochain-go/internal/logger"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

type BlockChain struct {
	ctx   context.Context
	Block []*Block
}

func NewBlockChain(ctx context.Context, tp tm.TimeProvider) *BlockChain {
	genesis := newGenesisBlock(tp.Now())
	blockChain := &BlockChain{
		ctx:   ctx,
		Block: []*Block{genesis},
	}

	return blockChain
}

func (bc *BlockChain) AddBlock(data string, tm tm.TimeProvider) {
	lastBlock := bc.Block[len(bc.Block)-1]
	addBlock := MineBlock(lastBlock, data, tm)

	bc.Block = append(bc.Block, addBlock)
}

func (bc *BlockChain) IsValidChain() bool {
	genesis := newGenesisBlock(bc.Block[0].Timestamp)
	if !reflect.DeepEqual(bc.Block[0], genesis) {
		return false
	}

	actualLastHash := genesis.Hash
	lastDifficulty := genesis.Difficulty
	for _, block := range bc.Block[1:] {
		if actualLastHash != block.LastHash {
			return false
		}
		nonce := block.Nonce
		difficulty := block.Difficulty
		validatedHash := cryptoHash(block.Timestamp.String(), strconv.Itoa(nonce), strconv.Itoa(difficulty), block.LastHash, block.Data)
		if block.Hash != validatedHash {
			return false
		}
		if math.Abs(float64(lastDifficulty-difficulty)) > 1 {
			return false
		}
		actualLastHash = block.Hash
		lastDifficulty = block.Difficulty
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
