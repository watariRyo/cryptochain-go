package block

import (
	"context"
	"encoding/json"
	"math"
	"reflect"
	"strconv"

	"github.com/watariRyo/cryptochain-go/internal/crypto"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	"github.com/watariRyo/cryptochain-go/internal/time"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
)

type BlockChain struct {
	ctx          context.Context
	block        []*model.Block
}

var _ repository.BlockChainInterface = (*BlockChain)(nil)

func NewBlockChain(ctx context.Context, tp tm.TimeProvider) *BlockChain {
	genesis := newGenesisBlock(tp.NowMicroString())
	blockChain := &BlockChain{
		ctx:   ctx,
		block: []*model.Block{genesis},
	}

	return blockChain
}

func (bc *BlockChain) AddBlock(data string, timeProvider tm.TimeProvider) {
	lastBlock := bc.block[len(bc.block)-1]
	addBlock := mineBlock(lastBlock, data, timeProvider)

	bc.block = append(bc.block, addBlock)
}

func (bc *BlockChain) GetBlock() []*model.Block {
	return bc.block
}

func (bc *BlockChain) IsValidChain() bool {
	genesis := newGenesisBlock(bc.block[0].Timestamp)
	if !reflect.DeepEqual(bc.block[0], genesis) {
		logger.Debugf(bc.ctx, "genesis")
		return false
	}

	actualLastHash := genesis.Hash
	lastDifficulty := genesis.Difficulty
	for _, block := range bc.block[1:] {
		if actualLastHash != block.LastHash {
			logger.Debugf(bc.ctx, "lastHash")
			return false
		}
		nonce := block.Nonce
		difficulty := block.Difficulty

		validatedHash := crypto.CryptoHash(block.Timestamp, strconv.Itoa(nonce), strconv.Itoa(difficulty), block.LastHash, block.Data)
		if block.Hash != validatedHash {
			logger.Debugf(bc.ctx, "hash")
			return false
		}
		if math.Abs(float64(lastDifficulty-difficulty)) > 1 {
			logger.Debugf(bc.ctx, "difficulty")
			return false
		}
		actualLastHash = block.Hash
		lastDifficulty = block.Difficulty
	}

	return true
}

func (bc *BlockChain) ReplaceChain(block []*model.Block, tm time.TimeProvider) {
	if len(block) <= len(bc.block) {
		logger.Warnf(bc.ctx, "The incoming chain must be longer.")
		return
	}

	checkChain := &BlockChain{
		ctx: bc.ctx,
		block: block,
	}

	if !checkChain.IsValidChain() {
		logger.Errorf(bc.ctx, "The incoming chain must be valid.")
		return
	}

	bc.block = block
}

func (bc *BlockChain) UnmarshalAndReplaceBlock(payload []byte, tm time.TimeProvider) {
	var payloadBlock []*model.Block
	if err := json.Unmarshal(payload, &payloadBlock); err != nil {
		logger.Errorf(bc.ctx, "Could not unmarshal block chain. %v", err)
	}
	subscribeChain := &BlockChain{
		ctx:   bc.ctx,
		block: payloadBlock,
	}
	bc.ReplaceChain(subscribeChain.block, tm)
}
