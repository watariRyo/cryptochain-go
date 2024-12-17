package usecase

import (
	"context"
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/watariRyo/cryptochain-go/internal/ec"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	mockRepository "github.com/watariRyo/cryptochain-go/web/domain/repository/mock"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
	"go.uber.org/mock/gomock"
)

func Test_GetTransactionPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrWallets := mockRepository.NewMockWalletsInterface(ctrl)
	mrWallets.EXPECT().GetTransactionPool().Times(1)

	uc := &UseCase{
		repo: &repository.AllRepository{Wallets: mrWallets},
	}

	uc.GetTransactionPool()
}

func Test_TransactExist(t *testing.T) {
	recipient := "test"
	amount := 50

	uc := testTransact(0, 1, amount, true, recipient, t)

	uc.Transact(&model.Transact{
		Amount:    amount,
		Recipient: recipient,
	})
}

func Test_TransactNotExist(t *testing.T) {
	recipient := "test"
	amount := 50

	uc := testTransact(1, 0, amount, false, recipient, t)

	uc.Transact(&model.Transact{
		Amount:    amount,
		Recipient: recipient,
	})
}

func testTransact(createTransactoinCnt, transactionUpdateCnt, amount int, isTransactionExist bool, recipient string, t *testing.T) *UseCase {
	ctrl := gomock.NewController(t)
	mrWallets := mockRepository.NewMockWalletsInterface(ctrl)
	dummyWallet := &model.Wallet{
		Curve:     ec.Secp256k1(),
		Balance:   1000,
		KeyPair:   &ecdsa.PrivateKey{},
		PublicKey: "publick-key",
	}
	mrWallets.EXPECT().GetWallet().Return(dummyWallet).Times(1)
	mrWallets.EXPECT().ExistingTransaction().Return(isTransactionExist).Times(1)
	mrWallets.EXPECT().GetTransaction().Times(1)
	mrWallets.EXPECT().SetTransaction(gomock.Any()).Times(1)
	mrWallets.EXPECT().GetTransactionPool().Times(1)

	mockTime := time.Date(2023, 12, 1, 12, 0, 0, 0, time.Local)
	mockTimeProvider := &MockTimeProvider{MockTime: mockTime}

	mrWallets.EXPECT().CreateTransaction(recipient, amount, mockTimeProvider).Times(createTransactoinCnt)
	mrWallets.EXPECT().TransactionUpdate(gomock.Any(), recipient, amount, mockTimeProvider).Times(transactionUpdateCnt)

	mrRedis := mockRepository.NewMockRedisClientInterface(ctrl)
	ctx := context.Background()
	mrRedis.EXPECT().Publish(ctx, (redis.TRANSACTION), gomock.Any()).MaxTimes(1)

	return &UseCase{
		ctx:          ctx,
		timeProvider: mockTimeProvider,
		repo:         &repository.AllRepository{Wallets: mrWallets, RedisClient: mrRedis},
	}
}
