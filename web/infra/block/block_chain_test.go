package block

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/watariRyo/cryptochain-go/internal/crypto"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

func TestChainStartWithGenesis(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockProvider := &MockTimeProvider{MockTime: mockTime}

	blockChain := NewBlockChain(context.Background(), mockProvider)

	genesis := newGenesisBlock(mockProvider.NowMicroString())

	if d := cmp.Diff(blockChain.block[0], genesis); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestAddNewBlockChain(t *testing.T) {
	want := `{ "foo": "bar" }`
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockProvider := &MockTimeProvider{MockTime: mockTime}

	blockChain := NewBlockChain(context.Background(), mockProvider)

	blockChain.AddBlock(want, mockProvider)
	got := blockChain.block[len(blockChain.block)-1].Data
	if got != want {
		t.Errorf("add block chain mismatch: got %v, want %v", got, want)
	}
	lastHashAddBlock := blockChain.block[len(blockChain.block)-1].LastHash
	hashLastBlock := blockChain.block[len(blockChain.block)-2].Hash
	if lastHashAddBlock != hashLastBlock {
		t.Errorf("add block chain hash mismatch: lastHashNowBlock %v, hashLastBlock %v", lastHashAddBlock, hashLastBlock)
	}
}

func TestValidChain(t *testing.T) {
	t.Run("when the chain does not start with the genesis block", func(t *testing.T) {
		mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
		mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

		blockChain := NewBlockChain(context.Background(), mockTimeProvider)
		blockChain.block[0].Data = `fake-genesis`
		isValidChain := blockChain.IsValidChain()
		if isValidChain {
			t.Errorf("When the chain does not start with the genesis block. isValidChain should be false.")
		}
	})

	t.Run("when the chain starts with the genesis block and has multiple blocks", func(t *testing.T) {
		realTimeProvider := &tm.RealTimeProvider{}

		t.Run("and a lastHash reference has changed", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			blockChain.block[2].LastHash = "broken-lastHash"

			isValidChain := blockChain.IsValidChain()
			if isValidChain {
				t.Errorf("When the lastHash reference has changed. isValidChain should be false.")
			}
		})

		t.Run("and the chain contains a block with an invalid field", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			blockChain.block[2].Data = "some-bad-and-evil-data"

			isValidChain := blockChain.IsValidChain()
			if isValidChain {
				t.Errorf("When the chain contains a block with an invalid field. isValidChain should be false.")
			}
		})

		t.Run("and the chain contains a block with a jumped difficulty", func(t *testing.T) {
			// returns false
			blockChain := NewBlockChain(context.Background(), realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			lastBlock := blockChain.block[len(blockChain.block)-1]
			lastHash := lastBlock.Hash

			timestamp := realTimeProvider.NowMicroString()
			nonce := 0
			data := `hoge`
			difficulty := lastBlock.Difficulty - 3

			hash := crypto.CryptoHash(timestamp, lastHash, strconv.Itoa(difficulty), strconv.Itoa(nonce), data)
			badBlock := newBlock(timestamp, lastHash, hash, data, nonce, difficulty)

			blockChain.block = append(blockChain.block, badBlock)

			if blockChain.IsValidChain() {
				t.Errorf("should be invalid block chain")
			}
		})

		t.Run("and the chain does not contain any invalid blocks", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			isValidChain := blockChain.IsValidChain()
			if !isValidChain {
				t.Errorf("When the chain does not contain any invalid blocks. isValidChain should be true.")
			}
		})
	})
}

func TestReplaceChain(t *testing.T) {
	t.Run("When the new chain is not longer. does not replace the chain.", func(t *testing.T) {
		realTimeProvider := &tm.RealTimeProvider{}

		blockChain := NewBlockChain(context.Background(), realTimeProvider)
		originalChain := blockChain.block
		newChain := NewBlockChain(context.Background(), realTimeProvider)

		newChain.block[0].Data = `{ "new" : "new-chain" }`

		blockChain.ReplaceChain(newChain.block, realTimeProvider)

		if blockChain.block[0].Data != originalChain[0].Data {
			t.Errorf("When the new chain is not longer. should not replace the chain")
		}
	})

	t.Run("When the new chain is longer.", func(t *testing.T) {
		mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
		mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
		realTimeProvider := &tm.RealTimeProvider{}
		t.Run("when the chain is invalid does not replace the chain.", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), mockTimeProvider)
			originalChain := blockChain.block
			newChain := NewBlockChain(context.Background(), realTimeProvider)

			newChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			newChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			newChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			newChain.block[2].Hash = "some-fake-hash"

			blockChain.ReplaceChain(newChain.block, realTimeProvider)

			if !reflect.DeepEqual(blockChain.block, originalChain) {
				t.Errorf("When the chain is invalid. should not replace the chain")
			}
		})

		t.Run("when the chain is valid replace the chain.", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), mockTimeProvider)
			newChain := NewBlockChain(context.Background(), realTimeProvider)

			newChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			newChain.AddBlock(`{"data": "Bears"}`, realTimeProvider)
			newChain.AddBlock(`{"data": "Battlestar Galactica"}`, realTimeProvider)

			blockChain.ReplaceChain(newChain.block, realTimeProvider)

			if !reflect.DeepEqual(blockChain.block[2], newChain.block[2]) {
				t.Errorf("Chain should be replaced.")
			}
		})

		t.Run("when the chain is valid replace the chain from byte array.", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background(), mockTimeProvider)
			newChain := `
[
	{
		"Timestamp": "2024-12-11T06:18:29.851186Z",
		"LastHash": "____",
		"Hash": "hash-one",
		"Difficulty": 3,
		"Nonce": 0,
		"Data": "{ \"one\": \"one\" }"
	},
	{
		"Timestamp": "2024-12-11T06:18:41.054109Z",
		"LastHash": "hash-one",
		"Hash": "011bb76eb7a896c2838f5330d543001bd60704a441b221805ab871ed54ed816d",
		"Difficulty": 2,
		"Nonce": 5,
		"Data": "hoge7"
	},
	{
		"Timestamp": "2024-12-11T06:18:46.045773Z",
		"LastHash": "011bb76eb7a896c2838f5330d543001bd60704a441b221805ab871ed54ed816d",
		"Hash": "17d5667d15847f6b8613c1b6fddb254446065054f1f06e693a97ef9b56fdbca1",
		"Difficulty": 1,
		"Nonce": 2,
		"Data": "hoge7"
	}
]`

			blockChain.UnmarshalAndReplaceBlock([]byte(newChain), realTimeProvider, nil)

			if blockChain.block[1].Hash != "011bb76eb7a896c2838f5330d543001bd60704a441b221805ab871ed54ed816d" {
				t.Errorf("Chain should be replaced.")
			}
			if blockChain.block[2].Hash != "17d5667d15847f6b8613c1b6fddb254446065054f1f06e693a97ef9b56fdbca1" {
				t.Errorf("Chain should be replaced.")
			}
		})
	})
}
