package wallets

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
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
}

func Test_CaluculateBalanceAddSomeOfAllOutputsToTheWalletBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	w1, _ := NewWallet()
	wallets := NewWallets(w1, nil)
	w2, _ := NewWallet()
	wallets2 := NewWallets(w2, nil)

	wallets.CreateTransaction(w.PublicKey, 50, nil, mockTimeProvider)
	transactionOne := *wallets.Transaction

	wallets2.CreateTransaction(w.PublicKey, 60, nil, mockTimeProvider)
	transactionTwo := *wallets2.Transaction

	blockChain := block.NewBlockChain(context.TODO(), mockTimeProvider)

	var transactions []*model.Transaction
	transactions = append(transactions, &transactionOne, &transactionTwo)
	transactionsByte, _ := json.Marshal(transactions)
	blockChain.AddBlock(string(transactionsByte), mockTimeProvider)

	got, err := wallets.CaluculateBalance(blockChain.GetBlock(), w.PublicKey)
	if err != nil {
		t.Errorf("calculate balance failed. err: %v", err)
	}
	want := block.STARTING_BALANCE + transactionOne.OutputMap[w.PublicKey] + transactionTwo.OutputMap[w.PublicKey]
	if got != want {
		t.Errorf("calculate balance failed to calculate. got: %d want: %d", got, want)
	}
}

func Test_OutputAmountOfRecentTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	wallets := NewWallets(w, nil)
	wallets.CreateTransaction("foo-address", 30, nil, mockTimeProvider)

	var transactions []*model.Transaction
	transactions = append(transactions, wallets.Transaction)

	blockChain := block.NewBlockChain(context.TODO(), mockTimeProvider)
	json, _ := json.Marshal(transactions)
	blockChain.AddBlock(string(json), mockTimeProvider)

	balance, err := wallets.CaluculateBalance(blockChain.GetBlock(), wallets.Wallet.PublicKey)
	if err != nil {
		t.Errorf("Something went wrong calculateBalance. %v", err)
	}
	if balance != wallets.Transaction.OutputMap[wallets.Wallet.PublicKey] {
		t.Errorf("CalculateBalance calc failed. got: %d want: %d", balance, wallets.Transaction.OutputMap[wallets.Wallet.PublicKey])
	}
}

func Test_OutputNextToAndAfterRecentTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	wallets := NewWallets(w, nil)
	wallets.CreateTransaction("later-foo-address", 60, nil, mockTimeProvider)

	var transactions []*model.Transaction

	recentTransaction := wallets.Transaction

	blockChain := block.NewBlockChain(context.TODO(), mockTimeProvider)

	wallets.NewRewardTransaction(mockTimeProvider)
	sameBlockTransaction := wallets.Transaction

	wallets.CreateTransaction(wallets.Wallet.PublicKey, 75, nil, mockTimeProvider)
	nextBlockTransaction := wallets.Transaction

	transactions = append(transactions, recentTransaction, sameBlockTransaction, nextBlockTransaction)
	json, _ := json.Marshal(transactions)
	blockChain.AddBlock(string(json), mockTimeProvider)

	// includes the output amounts in the returned balance
	got, err := wallets.CaluculateBalance(blockChain.GetBlock(), wallets.Wallet.PublicKey)
	if err != nil {
		t.Errorf("Something went wrong calculateBalance. %v", err)
	}
	want := recentTransaction.OutputMap[wallets.Wallet.PublicKey] + sameBlockTransaction.OutputMap[wallets.Wallet.PublicKey] + nextBlockTransaction.OutputMap[wallets.Wallet.PublicKey]
	if got != want {
		t.Errorf("CalculateBalance calc failed. got: %d want: %d", got, want)
	}
}
