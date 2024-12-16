package wallet

import (
	"testing"
)

func Test_SigningData(t *testing.T) {
	input := "data"
	wallet, err := newWallet()
	if err != nil {
		t.Errorf("could not create wallet. %v", err)
	}

	// 署名
	r, s, err := wallet.Sign(input)
	if err != nil {
		t.Errorf("Could not signed. %v", err)
	}

	t.Run("verifies a signature", func(t *testing.T) {
		if !wallet.VerifySignature(input, r, s) {
			t.Errorf("should be verified.")
		}
	})

	t.Run("does not verifies an invalid signature", func(t *testing.T) {
		dummy, _ := newWallet()

		r, s, _ = dummy.Sign(input)
		if wallet.VerifySignature(input, r, s) {
			t.Errorf("should be not verified.")
		}
	})
}
