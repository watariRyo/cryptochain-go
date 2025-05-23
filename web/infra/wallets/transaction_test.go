package wallets

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/watariRyo/cryptochain-go/internal/ec"
)

func TestTransactionOutputsAmountToRecipient(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50
	err := newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}
	if wallets.Transaction.OutputMap[dummyRecipient] != amount {
		t.Errorf("Invalid outputMap recipient key value. key: %s, map value: %d, expected: %d", dummyRecipient, wallets.Transaction.OutputMap[dummyRecipient], amount)
	}
}

func TestTransactionOutputRemainingBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50
	err := newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}
	if wallets.Transaction.OutputMap[wallets.Wallet.PublicKey] != wallets.Wallet.Balance-amount {
		t.Errorf("Invalid outputMap publicKey Balance. key: %s, map value: %d, expected: %d", wallets.Wallet.PublicKey, wallets.Transaction.OutputMap[wallets.Wallet.PublicKey], wallets.Wallet.Balance-amount)
	}
}

func TestTransactionInput(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)

	dummyRecipient := "dummy-reciepient"
	amount := 50
	err := newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}

	transaction := wallets.Transaction

	if transaction.Input.Amount != wallets.Wallet.Balance {
		t.Errorf("Should be sets the amount to the senderWallet balance")
	}

	if transaction.Input.Address != wallets.Wallet.PublicKey {
		t.Errorf("Should be sets the address to the senderWallet public key")
	}

	bytes, _ := json.Marshal(transaction.OutputMap)
	if !ec.VerifySignature(w.Curve, w.PublicKey, bytes, transaction.Input.Signature.R, transaction.Input.Signature.S) {
		t.Errorf("Should be true the VerifySignature")
	}
}

func TestValidTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50

	t.Run("Transaction is Valid", func(t *testing.T) {
		newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
		if !wallets.ValidTransaction(context.TODO()) {
			t.Errorf("validTransaction should be true.")
		}
	})

	t.Run("Transaction is invalid by outputMap", func(t *testing.T) {
		newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
		transaction := wallets.Transaction
		transaction.OutputMap[wallets.Wallet.PublicKey] = 9999
		if wallets.ValidTransaction(context.TODO()) {
			t.Errorf("validTransaction should be false by outputMap.")
		}
	})

	t.Run("Transaction is invalid by signature", func(t *testing.T) {
		newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
		transaction := wallets.Transaction
		dummyWallet, _ := NewWallet()

		dummyR, dummyS, _ := ec.Sign(dummyWallet.KeyPair, []byte("dummy"))
		transaction.Input.Signature.R = dummyR
		transaction.Input.Signature.S = dummyS

		if wallets.ValidTransaction(context.TODO()) {
			t.Errorf("validTransaction should be false by signature.")
		}
	})
}

func TestUpdate(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50

	newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	transaction := wallets.Transaction

	originalSignature := transaction.Input.Signature
	originalSenderOutput := transaction.OutputMap[wallets.Wallet.PublicKey]
	nextRecipient := "next-recipient"
	nextAmount := 50

	wallets.TransactionUpdate(wallets.Wallet, nextRecipient, nextAmount, mockTimeProvider)

	if transaction.OutputMap[nextRecipient] != nextAmount {
		t.Errorf("outputs the amount to the next recipient should be match next amount.")
	}

	if transaction.OutputMap[wallets.Wallet.PublicKey] != originalSenderOutput-nextAmount {
		t.Errorf("subtracts is missmatched.")
	}

	total := 0
	for _, value := range transaction.OutputMap {
		total += value
	}
	if transaction.Input.Amount != total {
		t.Errorf("could not maintains a total output that matches the input amount")
	}

	// cmp.DiffはPublicでないと比較できない
	if transaction.Input.Signature.R == originalSignature.R || transaction.Input.Signature.S == originalSignature.S {
		t.Errorf("could not re-signs the transaction")
	}
}

func TestUpdateAmountExceedsBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)

	dummyRecipient := "dummy-reciepient"
	amount := 50

	newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)

	// originalSignature := transaction.input.signature
	// originalSenderOutput := transaction.outputMap[senderWallet.PublicKey]
	// nextRecipient := "next-recipient"
	// nextAmount := 50

	err := wallets.TransactionUpdate(wallets.Wallet, "foo", 999999, mockTimeProvider)
	if err == nil {
		t.Errorf("amount exceeds balance. should be return error")
	}
}

func TestUpdateAddedRecipentAmount(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	w, _ := NewWallet()
	wallets := NewWallets(context.TODO(), w, nil)
	dummyRecipient := "dummy-reciepient"
	amount := 50

	newTransaction(wallets, dummyRecipient, amount, mockTimeProvider)
	transaction := wallets.Transaction

	// originalSignature := transaction.input.signature
	originalSenderOutput := transaction.OutputMap[wallets.Wallet.PublicKey]
	addedAmount := 80

	err := wallets.TransactionUpdate(wallets.Wallet, dummyRecipient, addedAmount, mockTimeProvider)
	if err != nil {
		t.Errorf("transaction update failed. %v", err)
	}

	if transaction.OutputMap[dummyRecipient] != amount+addedAmount {
		t.Errorf("could not added amount to the same recipent")
	}
	if transaction.OutputMap[wallets.Wallet.PublicKey] != originalSenderOutput-addedAmount {
		t.Errorf("could not substract original output to the same recipent")
	}
}

func TestNewRewardTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	minerWallet, _ := NewWallet()
	wallets := NewWallets(context.TODO(), minerWallet, nil)

	// create a transacion with the reward input
	wallets.NewRewardTransaction(mockTimeProvider)
	rewardTransaction := wallets.Transaction

	if rewardTransaction.Input.Address != REWARD_INPUT {
		t.Errorf("RewardInput is missmatched. got: %s want: %s", rewardTransaction.Input.Address, REWARD_INPUT)
	}

	// creates ones transaction for the minner with the MINING_REWARD
	if rewardTransaction.OutputMap[minerWallet.PublicKey] != MINING_REWARD {
		t.Errorf("Mining Reward is missmatched. got: %d want: %d", rewardTransaction.OutputMap[minerWallet.PublicKey], MINING_REWARD)
	}
}
