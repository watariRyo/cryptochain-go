package usecase

import (
	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func (u *UseCase) Transact(req *model.Transact) (map[uuid.UUID]*model.Transaction, error) {
	wallet := u.repo.Wallets.GetWallet()
	if u.repo.Wallets.ExistingTransaction() {
		err := u.repo.Wallets.TransactionUpdate(wallet, req.Recipient, req.Amount, u.timeProvider)
		if err != nil {
			return nil, err
		}
	} else {
		err := u.repo.Wallets.CreateTransaction(req.Recipient, req.Amount, u.timeProvider)
		if err != nil {
			return nil, err
		}
	}

	u.repo.Wallets.SetTransaction(u.repo.Wallets.GetTransaction())
	pool := u.repo.Wallets.GetTransactionPool()

	return pool, nil
}

func (u *UseCase) GetTransactionPool() map[uuid.UUID]*model.Transaction {
	return u.repo.Wallets.GetTransactionPool()
}
