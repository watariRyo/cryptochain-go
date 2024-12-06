package block

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestChainStartWithGenesis(t *testing.T) {
	blockChain := NewBlockChain(context.Background())
	genesis := newGenesisBlock()

	if d := cmp.Diff(blockChain.Block[0], genesis); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestAddNewBlockChain(t *testing.T) {
	want := `{ "foo": "bar" }`
	blockChain := NewBlockChain(context.Background())

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockProvider := &MockTimeProvider{MockTime: mockTime}

	blockChain.AddBlock(want, mockProvider.Now())
	got := blockChain.Block[len(blockChain.Block)-1].Data
	if got != want {
		t.Errorf("add block chain mismatch: got %v, want %v", got, want)
	}
	lastHashAddBlock := blockChain.Block[len(blockChain.Block)-1].LastHash
	hashLastBlock := blockChain.Block[len(blockChain.Block)-2].Hash
	if lastHashAddBlock != hashLastBlock {
		t.Errorf("add block chain hash mismatch: lastHashNowBlock %v, hashLastBlock %v", lastHashAddBlock, hashLastBlock)
	}
}

func TestValidChain(t *testing.T) {
	t.Run("when the chain does not start with the genesis block", func(t *testing.T) {
		blockChain := NewBlockChain(context.Background())
		blockChain.Block[0].Data = `fake-genesis`
		isValidChain := blockChain.IsValidChain()
		if isValidChain {
			t.Errorf("When the chain does not start with the genesis block. isValidChain should be false.")
		}
	})

	t.Run("when the chain starts with the genesis block and has multiple blocks", func(t *testing.T) {
		mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
		mockProvider := &MockTimeProvider{MockTime: mockTime}
		t.Run("and a lastHash reference has changed", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, mockProvider.Now())

			blockChain.Block[2].LastHash = "broken-lastHash"

			isValidChain := blockChain.IsValidChain()
			if isValidChain {
				t.Errorf("When the lastHash reference has changed. isValidChain should be false.")
			}
		})

		t.Run("and the chain contains a block with an invalid field", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, mockProvider.Now())

			blockChain.Block[2].Data = "some-bad-and-evil-data"

			isValidChain := blockChain.IsValidChain()
			if isValidChain {
				t.Errorf("When the chain contains a block with an invalid field. isValidChain should be false.")
			}
		})

		t.Run("and the chain does not contain any invalid blocks", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			blockChain.AddBlock(`{"data": "Battlestar Galactica"}`, mockProvider.Now())

			isValidChain := blockChain.IsValidChain()
			if !isValidChain {
				t.Errorf("When the chain does not contain any invalid blocks. isValidChain should be true.")
			}
		})
	})
}

func TestReplaceChain(t *testing.T) {
	t.Run("When the new chain is not longer. does not replace the chain.", func(t *testing.T) {
		blockChain := NewBlockChain(context.Background())
		originalChain := blockChain.Block
		newChain := NewBlockChain(context.Background())

		newChain.Block[0].Data = `{ "new" : "new-chain" }`

		blockChain.ReplaceChain(newChain)

		if blockChain.Block[0].Data != originalChain[0].Data {
			t.Errorf("When the new chain is not longer. should not replace the chain")
		}
	})

	t.Run("When the new chain is longer.", func(t *testing.T) {
		mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
		mockProvider := &MockTimeProvider{MockTime: mockTime}
		t.Run("when the chain is invalid does not replace the chain.", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background())
			originalChain := blockChain.Block
			newChain := NewBlockChain(context.Background())

			newChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			newChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			newChain.AddBlock(`{"data": "Battlestar Galactica"}`, mockProvider.Now())

			newChain.Block[2].Hash = "some-fake-hash"

			blockChain.ReplaceChain(newChain)

			if !reflect.DeepEqual(blockChain.Block, originalChain) {
				t.Errorf("When the chain is invalid. should not replace the chain")
			}
		})

		t.Run("when the chain is valid replace the chain.", func(t *testing.T) {
			blockChain := NewBlockChain(context.Background())
			newChain := NewBlockChain(context.Background())

			newChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			newChain.AddBlock(`{"data": "Bears"}`, mockProvider.Now())
			newChain.AddBlock(`{"data": "Battlestar Galactica"}`, mockProvider.Now())

			blockChain.ReplaceChain(newChain)

			if !reflect.DeepEqual(blockChain.Block[2], newChain.Block[2]) {
				t.Errorf("Chain should be replaced.")
			}
		})
	})
}
