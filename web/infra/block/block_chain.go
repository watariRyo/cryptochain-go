package block

import (
	"context"
	"encoding/json"
	"fmt"
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
	block []*model.Block
}

var _ repository.BlockChainInterface = (*BlockChain)(nil)

func NewBlockChain(tp tm.TimeProvider) *BlockChain {
	genesis := newGenesisBlock(tp.NowMicroString())
	blockChain := &BlockChain{
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
		return false
	}

	actualLastHash := genesis.Hash
	lastDifficulty := genesis.Difficulty
	for _, block := range bc.block[1:] {
		if actualLastHash != block.LastHash {
			return false
		}
		nonce := block.Nonce
		difficulty := block.Difficulty

		validatedHash := crypto.CryptoHash(block.Timestamp, strconv.Itoa(nonce), strconv.Itoa(difficulty), block.LastHash, block.Data)

		if block.Hash != validatedHash {
			return false
		}
		if !bc.isValidDifficulty(difficulty, block.Hash) {
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

func (bc *BlockChain) isValidDifficulty(difficulty int, blockHash string) bool {
	checkCount := difficulty
	binary := ""
	for _, char := range blockHash {
		value := crypto.CharToBinary(char)
		binary += fmt.Sprintf("%04b", value)
	}
	for idx, hashRune := range binary {
		if checkCount == idx {
			break
		}

		char := fmt.Sprint(string(hashRune))
		if char != "0" {
			return false
		}
	}

	return true
}

func (bc *BlockChain) ReplaceChain(ctx context.Context, block []*model.Block, tm time.TimeProvider, validTransactionDataFn func(ctx context.Context, param1 []*model.Block, param2 []*model.Block) bool) {
	if len(block) <= len(bc.block) {
		logger.Warnf(ctx, "The incoming chain must be longer.")
		return
	}

	checkChain := &BlockChain{
		block: block,
	}

	if !checkChain.IsValidChain() {
		logger.Errorf(ctx, "The incoming chain must be valid.")
		return
	}
	if validTransactionDataFn != nil {
		if !validTransactionDataFn(ctx, bc.block, block) {
			logger.Errorf(ctx, "Invalid transaction data. did not replace chain")
			return
		}
	}

	bc.block = block
}

func (bc *BlockChain) UnmarshalAndReplaceBlock(ctx context.Context, payload []byte, tm time.TimeProvider, fn func([]*model.Block) error, validTransactionDataFn func(ctx context.Context, param1 []*model.Block, param2 []*model.Block) bool) {
	var payloadBlock []*model.Block
	if err := json.Unmarshal(payload, &payloadBlock); err != nil {
		logger.Errorf(ctx, "Could not unmarshal block chain. %v", err)
		return
	}
	subscribeChain := &BlockChain{
		block: payloadBlock,
	}

	bc.ReplaceChain(ctx, subscribeChain.block, tm, validTransactionDataFn)

	if fn != nil {
		if err := fn(payloadBlock); err != nil {
			// ここに渡すのはclearTransactionのみ
			logger.Errorf(ctx, "Could not clear transaction. %v", err)
		}
	}
}
