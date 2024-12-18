package wallets

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
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

func Test_ValidTransactins(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(w, nil)

	var validWantTransactions []*model.Transaction

	for idx := range 10 {
		newTransaction(wallets, "any-recipient", 30, mockTimeProvider)
		transaction := wallets.Transaction

		if idx%3 == 0 {
			transaction.Input.Amount = 999999
		} else if idx%3 == 1 {
			r, s, _ := ec.Sign(w.KeyPair, []byte("foo"))
			transaction.Input.Signature = &model.Signature{
				R: r,
				S: s,
			}
		} else {
			validWantTransactions = append(validWantTransactions, transaction)
		}
		wallets.SetTransaction(transaction)
	}

	validGotTransaction := wallets.ValidTransactoins(context.TODO())

	if len(validWantTransactions) != len(validGotTransaction) {
		t.Errorf("could not get expected valid transactions, want: %d, got: %d", len(validWantTransactions), len(validGotTransaction))
	}

	for _, want := range validWantTransactions {
		isMatched := false
		wantUUID := want.Id
		for _, got := range validGotTransaction {
			if got.Id == wantUUID {
				isMatched = true
				break
			}
		}
		if !isMatched {
			t.Errorf("could not get expected valid transactions")
		}
	}
}
