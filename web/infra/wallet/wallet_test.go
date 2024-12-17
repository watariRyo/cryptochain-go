package wallet

import (
	"testing"
	"time"
)

func Test_SigningData(t *testing.T) {
	input := "data"
	wallet, err := NewWallet()
	if err != nil {
		t.Errorf("could not create wallet. %v", err)
	}

	// 署名
	r, s, err := wallet.Sign([]byte(input))
	if err != nil {
		t.Errorf("Could not signed. %v", err)
	}

	t.Run("verifies a signature", func(t *testing.T) {
		if !wallet.VerifySignature([]byte(input), r, s) {
			t.Errorf("should be verified.")
		}
	})

	t.Run("does not verifies an invalid signature", func(t *testing.T) {
		dummy, _ := NewWallet()

		r, s, _ = dummy.Sign([]byte(input))
		if wallet.VerifySignature([]byte(input), r, s) {
			t.Errorf("should be not verified.")
		}
	})
}

func Test_CreateTrunsaction(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}
	t.Run("the amount exceeds the balanace", func(t *testing.T) {
		wallet, _ := NewWallet()

		if _, err := wallet.CreateTransaction("foo-recipient", 999999, mockTimeProvider); err == nil {
			t.Errorf("Amount exceeds balance. Should return error.")
		}
	})

	t.Run("the amount is valid", func(t *testing.T) {
		wallet, _ := NewWallet()
		amount := 50
		recipient := "foo-recipient"
		transaction, err := wallet.CreateTransaction(recipient, amount, mockTimeProvider)
		if err != nil {
			t.Errorf("Something went wrong at CreateTransaction")
		}

		if transaction.input.address != wallet.PublicKey {
			t.Errorf("should match the transaction input with the wallet.")
		}
		// outputs the amount the recipient
		if transaction.outputMap[recipient] != amount {
			t.Errorf("should match outputs the amount the recipient")
		}
	})
}
