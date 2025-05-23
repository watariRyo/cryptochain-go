package wallets

import (
	"context"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
)

type Wallets struct {
	Wallet          *model.Wallet
	Transaction     *model.Transaction
	TransactionPool map[uuid.UUID]*model.Transaction
}

var _ repository.WalletsInterface = (*Wallets)(nil)

func NewWallets(ctx context.Context, wallet *model.Wallet, transaction *model.Transaction) *Wallets {
	return &Wallets{
		Wallet:          wallet,
		Transaction:     transaction,
		TransactionPool: make(map[uuid.UUID]*model.Transaction),
	}
}

func (w *Wallets) GetTransactionPool() map[uuid.UUID]*model.Transaction {
	return w.TransactionPool
}

func (w *Wallets) GetWallet() *model.Wallet {
	return w.Wallet
}

func (w *Wallets) GetTransaction() *model.Transaction {
	return w.Transaction
}
