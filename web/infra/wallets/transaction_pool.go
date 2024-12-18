package wallets

import (
	"context"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func (wtp *Wallets) SetTransaction(transaction *model.Transaction) {
	wtp.TransactionPool[transaction.Id] = transaction
}

func (wtp *Wallets) SetMap(transactoinPool map[uuid.UUID]*model.Transaction) {
	wtp.TransactionPool = transactoinPool
}

func (wtp *Wallets) SetTransactionPool(transactoinPool map[uuid.UUID]*model.Transaction) {
	wtp.TransactionPool = transactoinPool
}

func (wtp *Wallets) ExistingTransaction() bool {
	for _, transaction := range wtp.TransactionPool {
		if wtp.Transaction.Input.Address == transaction.Input.Address {
			return true
		}
	}
	return false
}

func (wtp *Wallets) ValidTransactoins(ctx context.Context) []*model.Transaction {
	var validTransactions []*model.Transaction
	for _, transaction := range wtp.TransactionPool {
		if wtp.validTransaction(ctx, transaction) {
			validTransactions = append(validTransactions, transaction)
		}
	}
	return validTransactions
}
