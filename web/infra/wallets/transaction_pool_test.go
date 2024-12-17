package wallets

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_TransactionPool(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50

	newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	transaction := wallets.Transaction
	wallets.SetTransaction(transaction)

	want := transaction
	got := wallets.TransactionPool[transaction.Id]

	if want.Id != got.Id {
		t.Errorf("differs id: (got: %v want: %v)\n", got.Id, want.Id)
	}
	if d := cmp.Diff(got.OutputMap, want.OutputMap); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func Test_ExistingTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	wallets := NewWallets(w, nil)
	wallets.CreateTransaction("hoge", 50, mockTimeProvider)

	if wallets.ExistingTransaction() {
		t.Errorf("should return false. not an existing transaction.")
	}

	wallets.SetTransaction(wallets.Transaction)
	if !wallets.ExistingTransaction() {
		t.Errorf("should return true. an existing transaction.")
	}
}
