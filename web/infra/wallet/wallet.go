package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
)

type Wallet struct {
	curve     elliptic.Curve
	Balance   int
	keyPair   *ecdsa.PrivateKey
	PublicKey string
}

var _ repository.WalletInterface = (*Wallet)(nil)

func NewWallet() (*Wallet, error) {
	curve := ec.Secp256k1()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		curve:     curve,
		Balance:   block.STARTING_BALANCE,
		keyPair:   privateKey,
		PublicKey: ec.PublicKeyToHexCompressed(privateKey.PublicKey.X, privateKey.Y),
	}, nil
}

func (w *Wallet) Sign(message []byte) (r, s *big.Int, err error) {
	return ec.Sign(w.keyPair, message)
}

func (w *Wallet) VerifySignature(message []byte, r, s *big.Int) bool {
	return ec.VerifySignature(w.curve, w.PublicKey, message, r, s)
}

func (w *Wallet) CreateTransaction(recipient string, amount int, tm tm.TimeProvider) (*Transaction, error) {
	if amount > w.Balance {
		return nil, fmt.Errorf("amount exceeds balance. amount:%d balance:%d", amount, w.Balance)
	}
	return NewTransaction(w, recipient, amount, tm)
}
