package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	mockRepository "github.com/watariRyo/cryptochain-go/web/domain/repository/mock"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
	"go.uber.org/mock/gomock"
)

func Test_MineTransactions(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mw := mockRepository.NewMockWalletsInterface(ctrl)
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	dummyValidTransactions := []*model.Transaction{}
	dummyValidTransactions = append(dummyValidTransactions, &model.Transaction{
		Id:        uuid.New(),
		OutputMap: make(map[string]int),
		Input: &model.Input{
			Address: "dummy1",
		},
	})
	dummyValidTransactions = append(dummyValidTransactions, &model.Transaction{
		Id:        uuid.New(),
		OutputMap: make(map[string]int),
		Input: &model.Input{
			Address: "dummy2",
		},
	})

	dummyRewardTransactoin := &model.Transaction{
		Id:        uuid.New(),
		OutputMap: make(map[string]int),
		Input: &model.Input{
			Address: "reward",
		},
	}

	mw.EXPECT().ValidTransactoins(ctx).Return(dummyValidTransactions).Times(1)
	mw.EXPECT().NewRewardTransaction(mockTimeProvider).Times(1)
	mw.EXPECT().GetTransaction().Return(dummyRewardTransactoin).Times(1)

	dummyValidTransactions = append(dummyValidTransactions, dummyRewardTransactoin)

	dummyTransactionsBytes, _ := json.Marshal(dummyValidTransactions)

	mb := mockRepository.NewMockBlockChainInterface(ctrl)

	mb.EXPECT().AddBlock(string(dummyTransactionsBytes), mockTimeProvider).Times(1)
	mb.EXPECT().GetBlock().Times(1)

	mr := mockRepository.NewMockRedisClientInterface(ctrl)
	mr.EXPECT().Publish(ctx, string(redis.BLOCKCHAIN), gomock.Any()).MaxTimes(1)

	mw.EXPECT().ClearBlockChainTransactions(gomock.Any()).Times(1)

	uc := &UseCase{
		ctx:          ctx,
		timeProvider: mockTimeProvider,
		repo:         &repository.AllRepository{BlockChain: mb, RedisClient: mr, Wallets: mw},
	}

	uc.MineTransactions()
}
