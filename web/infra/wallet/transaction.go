package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
)

type Transaction struct {
	id        uuid.UUID
	outputMap map[string]int
	input     *Input
}

type Input struct {
	timestamp time.Time
	amount    int
	address   string
	signature *signature
}

type signature struct {
	r *big.Int
	s *big.Int
	// v int
}

func NewTransaction(senderWallet *Wallet, recipient string, amount int, tm tm.TimeProvider) (*Transaction, error) {
	outputMap := createOutputMap(senderWallet, recipient, amount)
	input, err := createInput(tm, senderWallet, outputMap)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		id:        uuid.New(),
		outputMap: outputMap,
		input:     input,
	}, nil

}

func createInput(tm tm.TimeProvider, senderWallet *Wallet, outputMap map[string]int) (*Input, error) {
	signatureDate, err := json.Marshal(outputMap)
	if err != nil {
		return nil, err
	}

	r, s, err := senderWallet.Sign(signatureDate)
	if err != nil {
		return nil, err
	}

	return &Input{
		timestamp: tm.Now(),
		amount:    senderWallet.Balance,
		address:   senderWallet.PublicKey,
		signature: &signature{
			r: r,
			s: s,
		},
	}, nil
}

func createOutputMap(senderWallet *Wallet, recipient string, amount int) map[string]int {
	outputMap := make(map[string]int)

	outputMap[recipient] = amount
	outputMap[senderWallet.PublicKey] = senderWallet.Balance - amount

	return outputMap
}

func (t *Transaction) validTransaction(ctx context.Context) bool {
	total := 0
	for _, value := range t.outputMap {
		total += value
	}

	if t.input.amount != total {
		logger.Errorf(ctx, "Invalid transaction from %s", t.input.address)
		return false
	}

	bytes, err := json.Marshal(t.outputMap)
	if err != nil {
		logger.Errorf(ctx, "Invalid outputMap %v", t.outputMap)
		return false
	}

	if !ec.VerifySignature(ec.Secp256k1(), t.input.address, bytes, t.input.signature.r, t.input.signature.s) {
		return false
	}

	return true
}

func (t *Transaction) update(senderWallet *Wallet, recpient string, amount int, tm tm.TimeProvider) error {
	if amount > t.outputMap[senderWallet.PublicKey] {
		return fmt.Errorf("amount exceeds balance")
	}

	recipentAmount, ok := t.outputMap[recpient]
	if ok {
		t.outputMap[recpient] = recipentAmount + amount
		t.outputMap[senderWallet.PublicKey] -= (recipentAmount + amount)
	} else {
		t.outputMap[recpient] = amount
		t.outputMap[senderWallet.PublicKey] -= amount
	}

	newInput, err := createInput(tm, senderWallet, t.outputMap)
	if err != nil {
		return err
	}
	t.input = newInput

	return nil
}
