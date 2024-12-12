package block

import "time"

const MINE_RATE = time.Duration(1 * time.Second)
const INITIAL_DIFFICULTY = 3

type genesisBlock struct {
	timestamp  string
	lastHash   string
	hash       string
	difficulty int
	nonce      int
	data       string
}

func newGenesis(timesamp string) *genesisBlock {
	return &genesisBlock{
		timestamp:  timesamp,
		lastHash:   "____",
		hash:       "hash-one",
		difficulty: INITIAL_DIFFICULTY,
		nonce:      0,
		data:       `{ "one": "one" }`,
	}
}
