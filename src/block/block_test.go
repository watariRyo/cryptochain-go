package block

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateNewBlock(t *testing.T) {
	timestamp := time.Now()
	hash := "test_hash"
	lastHash := "test_last_hash"
	data := "test_data"
	nonce := 1
	difficulty := 1

	block := newBlock(timestamp, lastHash, hash, data, nonce, difficulty)

	if block.Timestamp != timestamp {
		t.Errorf("expected %s, got %s", timestamp, block.Timestamp)
	}
	if block.Hash != hash {
		t.Errorf("expected %s, got %s", hash, block.Hash)
	}
	if block.LastHash != lastHash {
		t.Errorf("expected %s, got %s", lastHash, block.LastHash)
	}
	if block.Data != data {
		t.Errorf("expected %s, got %s", data, block.Data)
	}
	if block.Nonce != nonce {
		t.Errorf("expected %d, got %d", nonce, block.Nonce)
	}
	if block.Difficulty != difficulty {
		t.Errorf("expected %d, got %d", difficulty, block.Difficulty)
	}
}

func TestGenesisBlock(t *testing.T) {
	got := *newGenesisBlock()
	genesis := newGenesis()

	want := Block{
		Timestamp:  genesis.timestamp,
		LastHash:   genesis.lastHash,
		Hash:       genesis.hash,
		Data:       genesis.data,
		Difficulty: genesis.difficulty,
		Nonce:      genesis.nonce,
	}

	if d := cmp.Diff(got, want, cmpopts.IgnoreFields(got, "Data")); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestMineBlock(t *testing.T) {
	lastBlock := newGenesisBlock()
	data := "mined data"

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockProvider := &MockTimeProvider{MockTime: mockTime}

	mineBlock := MineBlock(lastBlock, data, mockProvider)

	if mineBlock.LastHash != lastBlock.Hash {
		t.Errorf("LastHash and Hash are mismatched. lastHash = %v, hash %v", mineBlock.LastHash, lastBlock.Hash)
	}
	if mineBlock.Data != data {
		t.Errorf("mineBlock.Data and data are mismatched. mineBlock.Data = %v, data = %v", mineBlock.Data, data)
	}
	if mineBlock.Timestamp != mockTime {
		t.Errorf("mineBlock.Timestamp and timestamp are mismatched. mineBlock.Timestamp = %v, timestamp = %v", mineBlock.Timestamp, mockTime)
	}
	hashExpected := cryptoHash(mineBlock.Timestamp.String(), strconv.Itoa(mineBlock.Nonce), strconv.Itoa(mineBlock.Difficulty), lastBlock.Hash, data)
	if mineBlock.Hash != hashExpected {
		t.Errorf("mineBlock.Hash and expected are mismatched. mineBlock.Hash = %v, hash = %v", mineBlock.Hash, hashExpected)
	}
}

func TestMatchDifficultyCriteria(t *testing.T) {
	lastBlock := newGenesisBlock()
	data := "mined data"
	difficulty := 1

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockProvider := &MockTimeProvider{MockTime: mockTime}

	mineBlock := MineBlock(lastBlock, data, mockProvider)

	println(mineBlock.Hash)

	want := strings.Repeat("0", difficulty)
	got := mineBlock.Hash[:difficulty]

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}
