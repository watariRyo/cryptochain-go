package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

type Wallet struct {
	Curve     elliptic.Curve
	Balance   int
	KeyPair   *ecdsa.PrivateKey
	PublicKey string
}
