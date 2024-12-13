package usecase

import (
	"encoding/json"

	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func (u *UseCase) GetBlock() []*model.Block {
	return u.repo.BlockChain.GetBlock()
}

func (u *UseCase) Mine(payload string) error {

	u.repo.BlockChain.AddBlock(payload, u.timeProvider)
	broadcastChain, err := json.Marshal(u.repo.BlockChain.GetBlock())
	if err != nil {
		return err
	}

	go u.repo.RedisClient.Publish(u.ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	return nil
}