package wallets

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
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

func (ww *Wallets) CreateTransaction(recipient string, amount int, chain []*model.Block, tm tm.TimeProvider) error {
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
