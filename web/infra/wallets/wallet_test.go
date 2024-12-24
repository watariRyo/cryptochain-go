package wallets

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
)

func Test_SigningData(t *testing.T) {
	input := "data"
	w, err := NewWallet()
	if err != nil {
		t.Errorf("could not create wallet. %v", err)
	}

	// 署名
	r, s, err := ec.Sign(w.KeyPair, []byte(input))
	if err != nil {
		t.Errorf("Could not signed. %v", err)
	}

	t.Run("verifies a signature", func(t *testing.T) {
		if !ec.VerifySignature(w.Curve, w.PublicKey, []byte(input), r, s) {
			t.Errorf("should be verified.")
		}
	})

	t.Run("does not verifies an invalid signature", func(t *testing.T) {
		dw, _ := NewWallet()

		r, s, _ = ec.Sign(dw.KeyPair, []byte(input))
		if ec.VerifySignature(w.Curve, w.PublicKey, []byte(input), r, s) {
			t.Errorf("should be not verified.")
		}
	})
}

func Test_CreateTrunsaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	t.Run("the amount exceeds the balanace", func(t *testing.T) {
		w, _ := NewWallet()
		wallets := NewWallets(w, nil)

		if err := wallets.CreateTransaction("foo-recipient", 999999, nil, mockTimeProvider); err == nil {
			t.Errorf("Amount exceeds balance. Should return error.")
		}
	})

	t.Run("the amount is valid", func(t *testing.T) {
		w, _ := NewWallet()
		wallets := NewWallets(w, nil)
		amount := 50
		recipient := "foo-recipient"
		err := wallets.CreateTransaction(recipient, amount, nil, mockTimeProvider)
		if err != nil {
			t.Errorf("Something went wrong at CreateTransaction")
		}

		if wallets.Transaction.Input.Address != wallets.Wallet.PublicKey {
			t.Errorf("should match the transaction input with the wallet.")
		}
		// outputs the amount the recipient
		if wallets.Transaction.OutputMap[recipient] != amount {
			t.Errorf("should match outputs the amount the recipient")
		}
	})

	t.Run("a chain is passed", func(t *testing.T) {
		w, _ := NewWallet()
		wallets := NewWallets(w, nil)
		amount := 50
		recipient := "foo-recipient"

		// calls CalculateBalance
		blockChain := block.NewBlockChain(context.TODO(), mockTimeProvider)
		err := wallets.CreateTransaction(recipient, amount, blockChain.GetBlock(), mockTimeProvider)
		if err != nil {
			t.Errorf("Something went wrong at CreateTransaction")
		}
		if block.STARTING_BALANCE != wallets.Wallet.Balance {
			t.Errorf("Something went wrong at CreateTransaction")
		}
	})
}

func Test_CaluculateBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	wallets := NewWallets(w, nil)
	wallets.CreateTransaction("hoge", 50, nil, mockTimeProvider)

	blockChain := block.NewBlockChain(context.TODO(), mockTimeProvider)

	// return the STARTING_BALANCE
	gotStartingBalance, err := wallets.CaluculateBalance(blockChain.GetBlock(), wallets.Wallet.PublicKey)
	if err != nil {
		t.Errorf("calculate balance failed. err: %v", err)
	}
	if gotStartingBalance != block.STARTING_BALANCE {
		t.Errorf("failed to set genesis balanace")
	}

	wallets.CreateTransaction(wallets.Wallet.PublicKey, 50, nil, mockTimeProvider)
	transactionOne := wallets.Transaction
	wallets.CreateTransaction(wallets.Wallet.PublicKey, 60, nil, mockTimeProvider)
	transactionTwo := wallets.Transaction

	transactionOneBytes, _ := json.Marshal(transactionOne)
	transactionTwoBytes, _ := json.Marshal(transactionTwo)
	blockChain.AddBlock(string(transactionOneBytes), mockTimeProvider)
	blockChain.AddBlock(string(transactionTwoBytes), mockTimeProvider)

	// addds the sum of all outputs to the wallet balance
	got, err := wallets.CaluculateBalance(blockChain.GetBlock(), wallets.Wallet.PublicKey)
	if err != nil {
		t.Errorf("calculate balance failed. err: %v", err)
	}
	want := block.STARTING_BALANCE + transactionOne.OutputMap[wallets.Wallet.PublicKey] + transactionTwo.OutputMap[wallets.Wallet.PublicKey]
	if got != want {
		t.Errorf("calculate balance failed to calculate. got: %d want: %d", got, want)
	}
}
