package block

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

type Block struct {
	Timestamp  string
	LastHash   string
	Hash       string
	Difficulty int
	Nonce      int
	Data       string
}

func newBlock(timestamp string, lastHash, hash, data string, nonce, difficulty int) *Block {
	return &Block{
		Timestamp:  timestamp,
		LastHash:   lastHash,
		Hash:       hash,
		Difficulty: difficulty,
		Nonce:      nonce,
		Data:       data,
	}
}

func newGenesisBlock(timestamp string) *Block {
	gen := newGenesis(timestamp)
	return newBlock(gen.timestamp, gen.lastHash, gen.hash, gen.data, gen.nonce, gen.difficulty)
}

func mineBlock(lastBlock *Block, data string, tp tm.TimeProvider) *Block {
	nonce := 0

	difficulty := lastBlock.Difficulty
	var hash string
	var timestampStr string
	for {
		nonce++
		timestampStr = tp.NowMicroString()
		timestamp, _ := tm.MicroParse(timestampStr)
		difficulty = adjustDifficulty(lastBlock, timestamp)
		hash = cryptoHash(timestampStr, strconv.Itoa(nonce), strconv.Itoa(difficulty), lastBlock.Hash, data)
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
		Timestamp:  timestampStr,
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

	parseTime, _ := tm.MicroParse(originalBlock.Timestamp)

	difference := timestamp.Sub(parseTime)

	if difference > MINE_RATE {
		return difficulty - 1
	}

	return difficulty + 1
}
