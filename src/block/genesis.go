package block

import "time"

const INITIAL_DIFFICULTY = 3

type genesisBlock struct {
	timestamp  time.Time
	lastHash   string
	hash       string
	difficulty int
	nonce      int
	data       string
}

func newGenesis() *genesisBlock {
	return &genesisBlock{
		timestamp:  time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Now().Location()),
		lastHash:   "____",
		hash:       "hash-one",
		difficulty: INITIAL_DIFFICULTY,
		nonce:      0,
		data:       `{ "one": "one" }`,
	}
}
