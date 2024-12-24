package usecase

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func (u *UseCase) Transact(req *model.Transact) (map[uuid.UUID]*model.Transaction, error) {
	wallet := u.repo.Wallets.GetWallet()
	if u.repo.Wallets.ExistingTransaction() {
		err := u.repo.Wallets.TransactionUpdate(wallet, req.Recipient, req.Amount, u.timeProvider)
		if err != nil {
			return nil, err
		}
	} else {
		err := u.repo.Wallets.CreateTransaction(req.Recipient, req.Amount, u.GetBlock(), u.timeProvider)
		if err != nil {
			return nil, err
		}
	}

	transaction := u.repo.Wallets.GetTransaction()

	u.repo.Wallets.SetTransaction(transaction)

	broadcastTransaction, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	go u.repo.RedisClient.Publish(u.ctx, string(redis.TRANSACTION), string(broadcastTransaction))

	pool := u.repo.Wallets.GetTransactionPool()

	return pool, nil
}

func (u *UseCase) GetTransactionPool() map[uuid.UUID]*model.Transaction {
	return u.repo.Wallets.GetTransactionPool()
}
