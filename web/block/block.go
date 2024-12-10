package block

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

type Block struct {
	Timestamp  time.Time
	LastHash   string
	Hash       string
	Difficulty int
	Nonce      int
	Data       string
}

func newBlock(timestamp time.Time, lastHash, hash, data string, nonce, difficulty int) *Block {
	return &Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Difficulty: difficulty,
		Nonce:      nonce,
		Data:       data,
	}
}

func newGenesisBlock(timestamp time.Time) *Block {
	gen := newGenesis(timestamp)
	return newBlock(gen.timestamp, gen.lastHash, gen.hash, gen.data, gen.nonce, gen.difficulty)
}

func MineBlock(lastBlock *Block, data string, tp tm.TimeProvider) *Block {
	nonce := 0

	difficulty := lastBlock.Difficulty
	var hash string
	var timestamp time.Time

	for {
		nonce++
		timestamp = tp.Now()
		difficulty = adjustDifficulty(lastBlock, timestamp)
		hash = cryptoHash(timestamp.String(), strconv.Itoa(nonce), strconv.Itoa(difficulty), lastBlock.Hash, data)
		want := strings.Repeat("0", difficulty)

		binary := ""
		for _, char := range hash {
			value := charToBinary(char)
			binary += fmt.Sprintf("%04b", value)
		}

		if binary[:difficulty] == want {
			break
		}
	}

	return &Block{
		Timestamp:  timestamp,
		LastHash:   lastBlock.Hash,
		Difficulty: difficulty,
		Nonce:      nonce,
		Data:       data,
		Hash:       hash,
	}
}

func adjustDifficulty(originalBlock *Block, timestamp time.Time) int {
	difficulty := originalBlock.Difficulty

	if difficulty < 1 {
		return 1
	}

	difference := timestamp.Sub(originalBlock.Timestamp)

	if difference > MINE_RATE {
		return difficulty - 1
	}

	return difficulty + 1
}
