package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	mockRepository "github.com/watariRyo/cryptochain-go/web/domain/repository/mock"
)

func TestGetWalletInfo(t *testing.T) {
	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	ctrl := gomock.NewController(t)
	mrBlockChain := mockRepository.NewMockBlockChainInterface(ctrl)

	dummyReturn := []*model.Block{{
		Timestamp:  mockTimeProvider.NowMicroString(),
		LastHash:   "lash-hash",
		Hash:       "hash",
		Difficulty: 0,
		Nonce:      0,
		Data:       `[{ "data": "ok" }, { "data" : "good" }]`,
	}}
	mrBlockChain.EXPECT().GetBlock().Times(1).Return(dummyReturn)

	mrWallets := mockRepository.NewMockWalletsInterface(ctrl)
	mrWallets.EXPECT().GetWallet().Return(&model.Wallet{
		Balance:   1,
		PublicKey: "dummyKey",
	})
	mrWallets.EXPECT().CaluculateBalance(dummyReturn, "dummyKey").Times(1)

	uc := &UseCase{
		timeProvider: mockTimeProvider,
		repo:         &repository.AllRepository{BlockChain: mrBlockChain, Wallets: mrWallets},
	}

	uc.GetWalletInfo()
}
