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
}
