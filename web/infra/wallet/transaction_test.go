package wallet

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func Test_TransactionOutputsAmountToRecipient(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	wallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50
	transaction, err := NewTransaction(wallet, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}
	if transaction.outputMap[dummyRecipient] != amount {
		t.Errorf("Invalid outputMap recipient key value. key: %s, map value: %d, expected: %d", dummyRecipient, transaction.outputMap[dummyRecipient], amount)
	}
}

func Test_TransactionOutputRemainingBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	wallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50
	transaction, err := NewTransaction(wallet, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}
	if transaction.outputMap[wallet.PublicKey] != wallet.Balance-amount {
		t.Errorf("Invalid outputMap publicKey Balance. key: %s, map value: %d, expected: %d", wallet.PublicKey, transaction.outputMap[wallet.PublicKey], wallet.Balance-amount)
	}
}

func Test_TransactionInput(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	wallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50
	transaction, err := NewTransaction(wallet, dummyRecipient, amount, mockTimeProvider)
	if err != nil {
		t.Errorf("New transaction went wrong. %v", err)
	}

	if transaction.input.amount != wallet.Balance {
		t.Errorf("Should be sets the amount to the senderWallet balance")
	}

	if transaction.input.address != wallet.PublicKey {
		t.Errorf("Should be sets the address to the senderWallet public key")
	}

	bytes, _ := json.Marshal(transaction.outputMap)
	if !wallet.VerifySignature(bytes, transaction.input.signature.r, transaction.input.signature.s) {
		t.Errorf("Should be true the VerifySignature")
	}
}

func Test_ValidTransaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	senderWallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50

	t.Run("Transaction is Valid", func(t *testing.T) {
		transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)
		if !transaction.validTransaction(context.TODO()) {
			t.Errorf("validTransaction should be true.")
		}
	})

	t.Run("Transaction is invalid by outputMap", func(t *testing.T) {
		transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)
		transaction.outputMap[senderWallet.PublicKey] = 9999
		if transaction.validTransaction(context.TODO()) {
			t.Errorf("validTransaction should be false by outputMap.")
		}
	})

	t.Run("Transaction is invalid by signature", func(t *testing.T) {
		transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)
		dummyWallet, _ := NewWallet()
		dummyR, dummyS, _ := dummyWallet.Sign([]byte("dummy"))
		transaction.input.signature.r = dummyR
		transaction.input.signature.s = dummyS

		if transaction.validTransaction(context.TODO()) {
			t.Errorf("validTransaction should be false by signature.")
		}
	})
}

func Test_Update(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	senderWallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50

	transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)

	originalSignature := transaction.input.signature
	originalSenderOutput := transaction.outputMap[senderWallet.PublicKey]
	nextRecipient := "next-recipient"
	nextAmount := 50

	transaction.update(senderWallet, nextRecipient, nextAmount, mockTimeProvider)

	if transaction.outputMap[nextRecipient] != nextAmount {
		t.Errorf("outputs the amount to the next recipient should be match next amount.")
	}

	if transaction.outputMap[senderWallet.PublicKey] != originalSenderOutput-nextAmount {
		t.Errorf("subtracts is missmatched.")
	}

	total := 0
	for _, value := range transaction.outputMap {
		total += value
	}
	if transaction.input.amount != total {
		t.Errorf("could not maintains a total output that matches the input amount")
	}

	// cmp.DiffはPublicでないと比較できない
	if transaction.input.signature.r == originalSignature.r || transaction.input.signature.s == originalSignature.s {
		t.Errorf("could not re-signs the transaction")
	}
}

func Test_UpdateAmountExceedsBalance(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	senderWallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50

	transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)

	// originalSignature := transaction.input.signature
	// originalSenderOutput := transaction.outputMap[senderWallet.PublicKey]
	// nextRecipient := "next-recipient"
	// nextAmount := 50

	err := transaction.update(senderWallet, "foo", 999999, mockTimeProvider)
	if err == nil {
		t.Errorf("amount exceeds balance. should be return error")
	}
}

func Test_UpdateAddedRecipentAmount(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	senderWallet, _ := NewWallet()
	dummyRecipient := "dummy-reciepient"
	amount := 50

	transaction, _ := NewTransaction(senderWallet, dummyRecipient, amount, mockTimeProvider)

	// originalSignature := transaction.input.signature
	originalSenderOutput := transaction.outputMap[senderWallet.PublicKey]
	addedAmount := 80

	err := transaction.update(senderWallet, dummyRecipient, addedAmount, mockTimeProvider)
	if err != nil {
		t.Errorf("transaction update failed. %v", err)
	}

	if transaction.outputMap[dummyRecipient] != amount+addedAmount {
		t.Errorf("could not added amount to the same recipent")
	}
	if transaction.outputMap[senderWallet.PublicKey] != originalSenderOutput-(amount+addedAmount) {
		t.Errorf("could not substract original output to the same recipent")
	}
}
