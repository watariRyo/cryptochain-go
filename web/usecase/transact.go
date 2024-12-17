package usecase

import (
	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func (u *UseCase) Transact(req model.Transact) (map[uuid.UUID]*model.Transaction, error) {
	err := u.repo.Wallets.CreateTransaction(req.Recipient, req.Amount, u.timeProvider)
	if err != nil {
		return nil, err
	}

	pool := u.repo.Wallets.GetTransactionPool()

	return pool, nil
}
