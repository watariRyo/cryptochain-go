package wallets

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/internal/logger"
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

func (ww *Wallets) CreateTransaction(recipient string, amount int, chain []*model.Block, tm tm.TimeProvider) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	if len(chain) > 0 {
		balance, err := ww.CaluculateBalance(chain, ww.Wallet.PublicKey)
		if err != nil {
			return err
		}
		ww.Wallet.Balance = balance
	}

	if amount > ww.Wallet.Balance {
		return fmt.Errorf("amount exceeds balance. amount:%d balance:%d", amount, ww.Wallet.Balance)
	}
	return newTransaction(ww, recipient, amount, tm)
}

func (ww *Wallets) CaluculateBalance(chain []*model.Block, address string) (int, error) {
	hasConductedTransaction := false
	outputsTotal := 0

	for i := len(chain) - 1; i > 0; i-- {
		block := chain[i]

		var transactions []*model.Transaction
		if err := json.Unmarshal([]byte(block.Data), &transactions); err != nil {
			return 0, err
		}

		for _, tr := range transactions {
			if tr.Input.Address == address {
				hasConductedTransaction = true
			}
			addressOutput, ok := tr.OutputMap[address]
			if ok {
				outputsTotal = outputsTotal + addressOutput
			}
		}
		if hasConductedTransaction {
			break
		}
	}

	if hasConductedTransaction {
		return outputsTotal, nil
	} else {
		return block.STARTING_BALANCE + outputsTotal, nil
	}
}

func (ww *Wallets) ValidTransactionData(ctx context.Context, originalChain []*model.Block, chain []*model.Block) bool {
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		rewardTransactionCount := 0

		var transactions []*model.Transaction
		if err := json.Unmarshal([]byte(block.Data), &transactions); err != nil {
			return false
		}

		transactionMap := make(map[uuid.UUID]bool)

		for _, tr := range transactions {
			if tr.Input.Address == REWARD_INPUT {
				rewardTransactionCount += 1

				if rewardTransactionCount > 1 {
					logger.Errorf(ctx, "Miner reward exceed limit")
					return false
				}

				for _, value := range tr.OutputMap {
					if value != MINING_REWARD {
						logger.Errorf(ctx, "Miner reward amount is invalid")
						return false
					}
				}
			} else {
				if !ww.validTransaction(ctx, tr) {
					logger.Errorf(ctx, "Invalid Transaction Data")
					return false
				}

				trueBalance, err := ww.CaluculateBalance(chain[0:i], tr.Input.Address)
				if err != nil {
					logger.Errorf(ctx, "Failed to CalculateBalance")
					return false
				}
				if tr.Input.Amount != trueBalance {
					logger.Errorf(ctx, "Invalid input amount. expected: %d got: %d", trueBalance, tr.Input.Amount)
					return false
				}

				if _, ok := transactionMap[tr.Id]; ok {
					logger.Errorf(ctx, "An identical transaction appears more than once in the block")
					return false
				} else {
					transactionMap[tr.Id] = true
				}
			}
		}
	}
	return true
}
