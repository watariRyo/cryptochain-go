package wallets

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
)

func NewWallet() (*model.Wallet, error) {
	curve := ec.Secp256k1()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &model.Wallet{
		Curve:     curve,
		Balance:   block.STARTING_BALANCE,
		KeyPair:   privateKey,
		PublicKey: ec.PublicKeyToHexCompressed(privateKey.PublicKey.X, privateKey.Y),
	}, nil
}

func (ww *Wallets) CreateTransaction(recipient string, amount int, tm tm.TimeProvider) error {
	if amount > ww.Wallet.Balance {
		return fmt.Errorf("amount exceeds balance. amount:%d balance:%d", amount, ww.Wallet.Balance)
	}
	return newTransaction(ww, recipient, amount, tm)
}
