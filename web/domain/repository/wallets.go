package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

type WalletsInterface interface {
	CreateTransaction(recipient string, amount int, tm time.TimeProvider) error
	ValidTransaction(ctx context.Context) bool
	TransactionUpdate(senderWallet *model.Wallet, recpient string, amount int, tm time.TimeProvider) error
	GetTransactionPool() map[uuid.UUID]*model.Transaction
	GetWallet() *model.Wallet
	GetTransaction() *model.Transaction
	SetTransaction(transaction *model.Transaction)
	ExistingTransaction() bool
	SetMap(transactoinPool map[uuid.UUID]*model.Transaction)
	ValidTransactoins(ctx context.Context) []*model.Transaction
	ClearTransactionPool()
	ClearBlockChainTransactions(chain []*model.Block) error
}
