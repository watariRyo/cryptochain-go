package block

import (
	"context"
	"encoding/json"
	"math"
	"reflect"
	"strconv"

	"github.com/watariRyo/cryptochain-go/internal/logger"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

type BlockChain struct {
	Ctx          context.Context
	TimeProvider *tm.RealTimeProvider
	block        []*Block
}

type BlockChainInterface interface {
	AddBlock(data string)
	GetBlock() []*Block
	IsValidChain() bool
	ReplaceChain(chain *BlockChain)
	UnmarshalAndReplaceBlock(payload []byte)
}

var _ BlockChainInterface = (*BlockChain)(nil)

func NewBlockChain(ctx context.Context, tp tm.TimeProvider) *BlockChain {
	genesis := newGenesisBlock(tp.NowMicroString())
	blockChain := &BlockChain{
		Ctx:   ctx,
		block: []*Block{genesis},
	}

	return blockChain
}

func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.block[len(bc.block)-1]
	addBlock := mineBlock(lastBlock, data, bc.TimeProvider)

	bc.block = append(bc.block, addBlock)
}

func (bc *BlockChain) GetBlock() []*Block {
	return bc.block
}

func (bc *BlockChain) IsValidChain() bool {
	genesis := newGenesisBlock(bc.block[0].Timestamp)
	if !reflect.DeepEqual(bc.block[0], genesis) {
		logger.Debugf(bc.Ctx, "genesis")
		return false
	}

	actualLastHash := genesis.Hash
	lastDifficulty := genesis.Difficulty
	for _, block := range bc.block[1:] {
		if actualLastHash != block.LastHash {
			logger.Debugf(bc.Ctx, "lastHash")
			return false
		}
		nonce := block.Nonce
		difficulty := block.Difficulty

		validatedHash := cryptoHash(block.Timestamp, strconv.Itoa(nonce), strconv.Itoa(difficulty), block.LastHash, block.Data)
		if block.Hash != validatedHash {
			logger.Debugf(bc.Ctx, "hash")
			return false
		}
		if math.Abs(float64(lastDifficulty-difficulty)) > 1 {
			logger.Debugf(bc.Ctx, "difficulty")
			return false
		}
		actualLastHash = block.Hash
		lastDifficulty = block.Difficulty
	}

	return true
}

func (bc *BlockChain) ReplaceChain(chain *BlockChain) {
	if len(chain.block) <= len(bc.block) {
		logger.Warnf(bc.Ctx, "The incoming chain must be longer.")
		return
	}
	if !chain.IsValidChain() {
		logger.Errorf(bc.Ctx, "The incoming chain must be valid.")
		return
	}

	bc.block = chain.block
}

func (bc *BlockChain) UnmarshalAndReplaceBlock(payload []byte) {
	var payloadBlock []*Block
	if err := json.Unmarshal(payload, &payloadBlock); err != nil {
		logger.Errorf(bc.Ctx, "Could not unmarshal block chain. %v", err)
	}
	subscribeChain := &BlockChain{
		Ctx:   bc.Ctx,
		block: payloadBlock,
	}
	bc.ReplaceChain(subscribeChain)
}
