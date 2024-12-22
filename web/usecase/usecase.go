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
	ctx          context.Context
	timeProvider tm.TimeProvider
	repo         *repository.AllRepository
	configs      *configs.Config
}

type UseCaseInterface interface {
	GetBlock() []*model.Block
	Mine(payload string) error
	SyncWithRootState() error
	Transact(req *model.Transact) (map[uuid.UUID]*model.Transaction, error)
	GetTransactionPool() map[uuid.UUID]*model.Transaction
	MineTransactions() error
}

func NewUseCase(ctx context.Context, timeProvider tm.TimeProvider, repo *repository.AllRepository, configs *configs.Config) *UseCase {
	return &UseCase{
		ctx:          ctx,
		timeProvider: timeProvider,
		repo:         repo,
		configs:      configs,
	}
}
