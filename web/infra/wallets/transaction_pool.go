package wallets

import (
	"context"
	"encoding/json"

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

func (wtp *Wallets) ClearTransactionPool() {
	for u := range wtp.TransactionPool {
		delete(wtp.TransactionPool, u)
	}
}

func (wtp *Wallets) ClearBlockChainTransactions(chain []*model.Block) error {
	chainLength := len(chain)
	for i := 1; i < chainLength; i++ {
		block := chain[i]
		var transaction model.Transaction
		if err := json.Unmarshal([]byte(block.Data), &transaction); err != nil {
			var transactions []*model.Transaction
			if err := json.Unmarshal([]byte(block.Data), &transactions); err != nil {
				return err
			}
			for _, tr := range transactions {
				wtp.clearBlockChainTransaction(tr.Id)
			}
		}
		wtp.clearBlockChainTransaction(transaction.Id)
	}
	return nil
}

func (wtp *Wallets) clearBlockChainTransaction(id uuid.UUID) {
	_, ok := wtp.TransactionPool[id]
	if ok {
		delete(wtp.TransactionPool, id)
	}
}
