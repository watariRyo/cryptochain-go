package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	mockRepository "github.com/watariRyo/cryptochain-go/web/domain/repository/mock"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

func TestMine(t *testing.T) {
	addBlockParam := "addBlock"
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	ctrl := gomock.NewController(t)
	mrBlockChain := mockRepository.NewMockBlockChainInterface(ctrl)
	mrBlockChain.EXPECT().AddBlock(addBlockParam, mockTimeProvider).Times(1)

	dummyReturn := []*model.Block{{
		Timestamp:  mockTimeProvider.NowMicroString(),
		LastHash:   "lash-hash",
		Hash:       "hash",
		Difficulty: 0,
		Nonce:      0,
		Data:       `[{ "data": "ok" }, { "data" : "good" }]`,
	}}
	mrBlockChain.EXPECT().GetBlock().Times(1).Return(dummyReturn)

	mrRedis := mockRepository.NewMockRedisClientInterface(ctrl)
	ctx := context.Background()
	dummyParam, _ := json.Marshal(dummyReturn)
	// goroutineで呼ばれる前にテストが終わるためMaxTimes指定
	mrRedis.EXPECT().Publish(ctx, string(redis.BLOCKCHAIN), string(dummyParam)).MaxTimes(1)

	uc := &UseCase{
		timeProvider: mockTimeProvider,
		repo:         &repository.AllRepository{BlockChain: mrBlockChain, RedisClient: mrRedis},
	}

	uc.Mine(ctx, addBlockParam)
}

func TestGetBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrBlockChain := mockRepository.NewMockBlockChainInterface(ctrl)
	mrBlockChain.EXPECT().GetBlock().Times(1)

	uc := &UseCase{
		repo: &repository.AllRepository{BlockChain: mrBlockChain},
	}

	uc.GetBlock()
}
