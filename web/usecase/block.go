package usecase

import (
	"context"
	"encoding/json"

	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func (u *UseCase) GetBlock() []*model.Block {
	return u.repo.BlockChain.GetBlock()
}

func (u *UseCase) Mine(ctx context.Context, payload string) error {

	u.repo.BlockChain.AddBlock(payload, u.timeProvider)
	broadcastChain, err := json.Marshal(u.repo.BlockChain.GetBlock())
	if err != nil {
		return err
	}

	go u.repo.RedisClient.Publish(ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	return nil
}
