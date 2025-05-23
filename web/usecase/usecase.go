package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/watariRyo/cryptochain-go/configs"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
)

type UseCase struct {
	timeProvider tm.TimeProvider
	repo         *repository.AllRepository
	configs      *configs.Config
}

type UseCaseInterface interface {
	GetBlock() []*model.Block
	Mine(ctx context.Context, payload string) error
	SyncWithRootState(ctx context.Context) error
	Transact(ctx context.Context, req *model.Transact) (map[uuid.UUID]*model.Transaction, error)
	GetTransactionPool() map[uuid.UUID]*model.Transaction
	MineTransactions(ctx context.Context) error
	GetWalletInfo() (*model.WalletInfo, error)
}

func NewUseCase(timeProvider tm.TimeProvider, repo *repository.AllRepository, configs *configs.Config) *UseCase {
	return &UseCase{
		timeProvider: timeProvider,
		repo:         repo,
		configs:      configs,
	}
}
