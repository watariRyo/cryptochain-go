package wallets

import (
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func (wtp *Wallets) SetTransaction(transaction *model.Transaction) {
	wtp.TransactionPool[transaction.Id] = transaction
}

func (wtp *Wallets) ExistingTransaction() bool {
	for _, transaction := range wtp.TransactionPool {
		if wtp.Transaction.Input.Address == transaction.Input.Address {
			return true
		}
	}
	return false
}
