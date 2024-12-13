package model

import "context"

type BlockChain struct {
	Ctx   context.Context
	Block []*Block
}