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

func (ww *Wallets) CreateTransaction(recipient string, amount int, tm tm.TimeProvider) error {
	if amount > ww.Wallet.Balance {
		return fmt.Errorf("amount exceeds balance. amount:%d balance:%d", amount, ww.Wallet.Balance)
	}
	return newTransaction(ww, recipient, amount, tm)
}

func (ww *Wallets) CaluculateBalance(chain []*model.Block, address string) (int, error) {
	outputsTotal := 0

	for i := 1; i < len(chain); i++ {
		block := chain[i]

		var transaction model.Transaction
		if err := json.Unmarshal([]byte(block.Data), &transaction); err != nil {
			var transactions []*model.Transaction
			if err := json.Unmarshal([]byte(block.Data), &transactions); err != nil {
				return 0, err
			}
			for _, tr := range transactions {
				addressOutput, ok := tr.OutputMap[address]
				if ok {
					outputsTotal = outputsTotal + addressOutput
				}
			}
		}

		addressOutput, ok := transaction.OutputMap[address]
		if ok {
			outputsTotal = outputsTotal + addressOutput
		}
	}

	return block.STARTING_BALANCE + outputsTotal, nil
}
