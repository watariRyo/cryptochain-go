package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
)

func (u *UseCase) SyncWithRootState(ctx context.Context) error {
	// Sync Chain
	err := u.syncChain(ctx)
	if err != nil {
		return err
	}

	// Sync Transaction
	err = u.syncTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) syncChain(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blocks", u.configs.Host), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	payload, err := u.syncRequest(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	validTransactionDataFn := u.repo.Wallets.ValidTransactionData
	u.repo.BlockChain.UnmarshalAndReplaceBlock(ctx, payload, u.timeProvider, nil, validTransactionDataFn)

	return nil
}

func (u *UseCase) syncTransaction(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/transaction-pool-map", u.configs.Host), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	payload, err := u.syncRequest(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var payloadTransaction map[uuid.UUID]*model.Transaction
	if err := json.Unmarshal(payload, &payloadTransaction); err != nil {
		logger.Errorf(ctx, "Could not unmarshal transaction. %v", err)
	}

	u.repo.Wallets.SetMap(payloadTransaction)

	return nil
}

func (u *UseCase) syncRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Request Error. StatusCode: %d", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}
