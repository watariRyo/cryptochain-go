package wallets

import (
	"testing"
	"time"

	"github.com/watariRyo/cryptochain-go/internal/ec"
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

		if err := wallets.CreateTransaction("foo-recipient", 999999, mockTimeProvider); err == nil {
			t.Errorf("Amount exceeds balance. Should return error.")
		}
	})

	t.Run("the amount is valid", func(t *testing.T) {
		w, _ := NewWallet()
		wallets := NewWallets(w, nil)
		amount := 50
		recipient := "foo-recipient"
		err := wallets.CreateTransaction(recipient, amount, mockTimeProvider)
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
}
