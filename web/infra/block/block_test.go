package block

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func TestCreateNewBlock(t *testing.T) {
	timestamp := tm.MicroParseString(time.Now())
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
	timestamp := tm.MicroParseString(time.Now())
	got := *newGenesisBlock(timestamp)
	genesis := newGenesis(timestamp)

	want := model.Block{
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
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	lastBlock := newGenesisBlock(mockTimeProvider.NowMicroString())
	data := "mined data"

	mineBlock := mineBlock(lastBlock, data, mockTimeProvider)

	if mineBlock.LastHash != lastBlock.Hash {
		t.Errorf("LastHash and Hash are mismatched. lastHash = %v, hash %v", mineBlock.LastHash, lastBlock.Hash)
	}
	if mineBlock.Data != data {
		t.Errorf("mineBlock.Data and data are mismatched. mineBlock.Data = %v, data = %v", mineBlock.Data, data)
	}
	if mineBlock.Timestamp != mockTimeProvider.NowMicroString() {
		t.Errorf("mineBlock.Timestamp and timestamp are mismatched. mineBlock.Timestamp = %v, timestamp = %v", mineBlock.Timestamp, mockTimeProvider.NowMicroString())
	}
	hashExpected := cryptoHash(mineBlock.Timestamp, strconv.Itoa(mineBlock.Nonce), strconv.Itoa(mineBlock.Difficulty), lastBlock.Hash, data)
	if mineBlock.Hash != hashExpected {
		t.Errorf("mineBlock.Hash and expected are mismatched. mineBlock.Hash = %v, hash = %v", mineBlock.Hash, hashExpected)
	}
}

func TestMatchDifficultyCriteria(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	lastBlock := newGenesisBlock(mockTimeProvider.NowMicroString())
	data := "mined data"
	difficulty := 1

	mineBlock := mineBlock(lastBlock, data, mockTimeProvider)

	binary := ""
	for _, char := range mineBlock.Hash {
		value := charToBinary(char)
		binary += fmt.Sprintf("%04b", value)
	}

	want := strings.Repeat("0", difficulty)
	got := binary[:difficulty]

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestAdjustDifficulty(t *testing.T) {
	data := "mined data"

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	lastBlock := newGenesisBlock(mockTimeProvider.NowMicroString())
	lastBlock.Difficulty = 3

	mineBlock := mineBlock(lastBlock, data, mockTimeProvider)

	mineTimestamp, _ := tm.MicroParse(mineBlock.Timestamp)

	raiseTimestamp := mineTimestamp.Add(time.Duration(-1 * time.Second))

	// raise the difficulty
	newDifficulty := adjustDifficulty(mineBlock, raiseTimestamp)
	if newDifficulty != mineBlock.Difficulty+1 {
		t.Errorf("Expected difficulty to be raised to %d, got %d", newDifficulty, mineBlock.Difficulty+1)
	}
	// lowers the difficulty
	mineTimestamp, _ = tm.MicroParse(mineBlock.Timestamp)

	lowersTimestamp := mineTimestamp.Add(time.Duration(2 * time.Second))

	newDifficulty = adjustDifficulty(mineBlock, lowersTimestamp)
	if newDifficulty != mineBlock.Difficulty-1 {
		t.Errorf("Expected difficulty to be lowers to %d, got %d", newDifficulty, mineBlock.Difficulty-1)
	}

	// adjust the difficulty in mineBlock
	possibleResults := []int{lastBlock.Difficulty + 1, lastBlock.Difficulty - 1}

	if !slices.Contains(possibleResults, mineBlock.Difficulty) {
		t.Errorf("mineBlock.Difficulty should be adjusted.")
	}
}

func TestAdjustDifficultyLowerLimit(t *testing.T) {
	data := "mined data"

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	lastBlock := newGenesisBlock(mockTimeProvider.NowMicroString())
	lastBlock.Difficulty = 1

	mineBlock := mineBlock(lastBlock, data, mockTimeProvider)
	mineBlock.Difficulty = -1

	got := adjustDifficulty(mineBlock, mockTimeProvider.Now())
	if got != 1 {
		t.Errorf("difficulty should be limited to 1, got %v", got)
	}
}
