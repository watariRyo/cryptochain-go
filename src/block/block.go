package block

import (
	"strconv"
	"strings"
	"time"

	tm "github.com/watariRyo/cryptochain-go/src/time"
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

func newGenesisBlock() *Block {
	gen := newGenesis()
	return newBlock(gen.timestamp, gen.lastHash, gen.hash, gen.data, gen.nonce, gen.difficulty)
}

func MineBlock(lastBlock *Block, data string, tp tm.TimeProvider) *Block {
	difficulty := lastBlock.Difficulty
	nonce := 0
	want := strings.Repeat("0", difficulty)

	var hash string
	var timestamp time.Time

	for {
		nonce++
		timestamp = tp.Now()
		hash = cryptoHash(timestamp.String(), strconv.Itoa(nonce), strconv.Itoa(difficulty), lastBlock.Hash, data)
		if hash[:difficulty] == want {
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
