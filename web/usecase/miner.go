package usecase

import (
	"context"
	"encoding/json"

	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func (u *UseCase) MineTransactions(ctx context.Context) error {
	validTransactions := u.repo.Wallets.ValidTransactoins(ctx)

	err := u.repo.Wallets.NewRewardTransaction(u.timeProvider)
	if err != nil {
		return err
	}
	validTransactions = append(validTransactions, u.repo.Wallets.GetTransaction())

	validTransactionBytes, err := json.Marshal(validTransactions)
	if err != nil {
		return err
	}
	u.repo.BlockChain.AddBlock(string(validTransactionBytes), u.timeProvider)

	chain := u.repo.BlockChain.GetBlock()
	broadcastChain, err := json.Marshal(chain)
	if err != nil {
		return err
	}

	go u.repo.RedisClient.Publish(ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	u.repo.Wallets.ClearBlockChainTransactions(chain)

	return nil
}
