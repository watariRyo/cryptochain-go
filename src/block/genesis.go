package block

import "time"

type genesisBlock struct {
	timestamp time.Time
	lastHash  string
	hash      string
	data      string
}

func newGenesis() *genesisBlock {
	return &genesisBlock{
		timestamp: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Now().Location()),
		lastHash:  "____",
		hash:      "hash-one",
		data:      `{ "one": "one" }`,
	}
}
