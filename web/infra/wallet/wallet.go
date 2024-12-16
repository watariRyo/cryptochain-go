package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
)

type Wallet struct {
	curve     elliptic.Curve
	balance   int64
	keyPair   *ecdsa.PrivateKey
	publicKey string
}

func newWallet() (*Wallet, error) {
	curve := ec.Secp256k1()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		curve:     curve,
		balance:   block.STARTING_BALANCE,
		keyPair:   privateKey,
		publicKey: ec.PublicKeyToHexCompressed(privateKey.PublicKey.X, privateKey.Y),
	}, nil
}

func (w *Wallet) Sign(message string) (r, s *big.Int, err error) {
	hash := sha256.Sum256([]byte(message))
	// 署名作成
	r, s, err = ecdsa.Sign(rand.Reader, w.keyPair, hash[:])
	return r, s, err
}

func (w *Wallet) VerifySignature(message string, r, s *big.Int) bool {
	hash := sha256.Sum256([]byte(message))

	x, y, err := ec.DecompressHexPublicKey(w.curve, w.publicKey)
	if err != nil {
		return false
	}

	publicKey := &ecdsa.PublicKey{
		Curve: w.curve,
		X:     x,
		Y:     y,
	}

	// 署名検証
	return ecdsa.Verify(publicKey, hash[:], r, s)
}
