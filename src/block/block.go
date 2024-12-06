package block

import (
	"time"
)

type Block struct {
	Timestamp time.Time
	LastHash  string
	Hash      string
	Data      string
}

func newBlock(timestamp time.Time, lastHash, hash string, data string) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}

func newGenesisBlock() *Block {
	gen := newGenesis()
	return newBlock(gen.timestamp, gen.lastHash, gen.hash, gen.data)
}

func MineBlock(lastBlock *Block, data string, timestamp time.Time) *Block {
	return &Block{
		Timestamp: timestamp,
		LastHash:  lastBlock.Hash,
		Data:      data,
		Hash:      cryptoHash(timestamp.String(), lastBlock.Hash, data),
	}
}
