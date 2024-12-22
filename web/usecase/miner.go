package usecase

import (
	"encoding/json"

	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func (u *UseCase) MineTransactions() error {
	validTransactions := u.repo.Wallets.ValidTransactoins(u.ctx)

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

	broadcastChain, err := json.Marshal(u.repo.BlockChain.GetBlock())
	if err != nil {
		return err
	}

	go u.repo.RedisClient.Publish(u.ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	u.repo.Wallets.ClearTransactionPool()

	return nil
}
