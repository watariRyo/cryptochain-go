package wallets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

const MINING_REWARD int = 50
const REWARD_INPUT = "*authorized-reward*"

func newTransaction(senderWallet *Wallets, recipient string, amount int, tm tm.TimeProvider) error {
	outputMap := createOutputMap(senderWallet.Wallet, recipient, amount)
	input, err := createInput(tm, senderWallet.Wallet, outputMap)
	if err != nil {
		return err
	}

	transactoin := &model.Transaction{
		Id:        uuid.New(),
		OutputMap: outputMap,
		Input:     input,
	}

	senderWallet.Transaction = transactoin
	return nil
}

func createInput(tm tm.TimeProvider, senderWallet *model.Wallet, outputMap map[string]int) (*model.Input, error) {
	signatureDate, err := json.Marshal(outputMap)
	if err != nil {
		return nil, err
	}

	r, s, err := ec.Sign(senderWallet.KeyPair, signatureDate)
	if err != nil {
		return nil, err
	}

	return &model.Input{
		Timestamp: tm.Now(),
		Amount:    senderWallet.Balance,
		Address:   senderWallet.PublicKey,
		Signature: &model.Signature{
			R: r,
			S: s,
		},
	}, nil
}

func createOutputMap(senderWallet *model.Wallet, recipient string, amount int) map[string]int {
	outputMap := make(map[string]int)

	outputMap[recipient] = amount
	outputMap[senderWallet.PublicKey] = senderWallet.Balance - amount

	return outputMap
}

func (wt *Wallets) NewRewardTransaction(tm tm.TimeProvider) error {
	outputMap := make(map[string]int)
	outputMap[wt.Wallet.PublicKey] = MINING_REWARD
	input := &model.Input{
		Address: REWARD_INPUT,
	}

	transactoin := &model.Transaction{
		Id:        uuid.New(),
		OutputMap: outputMap,
		Input:     input,
	}

	wt.Transaction = transactoin
	return nil
}

func (wt *Wallets) ValidTransaction(ctx context.Context) bool {
	return wt.validTransaction(ctx, wt.Transaction)
}

func (wt *Wallets) validTransaction(ctx context.Context, transaction *model.Transaction) bool {
	total := 0
	for _, value := range transaction.OutputMap {
		total += value
	}

	if transaction.Input.Amount != total {
		logger.Errorf(ctx, "Invalid transaction from %s. amount: %d total: %d", transaction.Input.Address, transaction.Input.Amount, total)
		return false
	}

	bytes, err := json.Marshal(transaction.OutputMap)
	if err != nil {
		logger.Errorf(ctx, "Invalid outputMap %v", transaction.OutputMap)
		return false
	}

	if !ec.VerifySignature(ec.Secp256k1(), transaction.Input.Address, bytes, transaction.Input.Signature.R, transaction.Input.Signature.S) {
		return false
	}

	return true
}

func (wt *Wallets) TransactionUpdate(senderWallet *model.Wallet, recpient string, amount int, tm tm.TimeProvider) error {
	if amount > wt.Transaction.OutputMap[senderWallet.PublicKey] {
		return fmt.Errorf("amount exceeds balance")
	}

	recipentAmount, ok := wt.Transaction.OutputMap[recpient]
	if ok {
		wt.Transaction.OutputMap[recpient] = recipentAmount + amount
	} else {
		wt.Transaction.OutputMap[recpient] = amount
	}

	wt.Transaction.OutputMap[senderWallet.PublicKey] -= amount

	newInput, err := createInput(tm, senderWallet, wt.Transaction.OutputMap)
	if err != nil {
		return err
	}
	wt.Transaction.Input = newInput

	return nil
}
