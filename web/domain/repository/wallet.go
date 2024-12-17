package repository

import "math/big"

type WalletInterface interface {
	Sign(message []byte) (r, s *big.Int, err error)
	VerifySignature(message []byte, r, s *big.Int) bool
}
