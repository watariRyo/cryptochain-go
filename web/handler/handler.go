package handler

import (
	"context"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/usecase"
)

type Handler struct {
	ctx     context.Context
	configs *configs.Config
	usecase usecase.UseCaseInterface
}

type HandlerInterface interface {
	GetBlocks(w http.ResponseWriter, r *http.Request)
	Mine(w http.ResponseWriter, r *http.Request)
	Transact(w http.ResponseWriter, r *http.Request)
	GetTransactionPool(w http.ResponseWriter, r *http.Request)
	GetMineTransactions(w http.ResponseWriter, r *http.Request)
}

var _ HandlerInterface = (*Handler)(nil)

func NewHandler(ctx context.Context, usecase *usecase.UseCase, configs *configs.Config) *Handler {
	return &Handler{
		ctx:     ctx,
		usecase: usecase,
		configs: configs,
	}
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusOK, handler.usecase.GetBlock())
}

func (handler *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	var requestPayload model.Payload

	handler.readJSON(w, r, &requestPayload)

	err := handler.usecase.Mine(requestPayload.Data)
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	handler.GetBlocks(w, r)
}

func (handler *Handler) Transact(w http.ResponseWriter, r *http.Request) {
	var requestPayload model.Transact

	handler.readJSON(w, r, &requestPayload)

	pool, err := handler.usecase.Transact(&requestPayload)
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	handler.writeJSON(w, http.StatusOK, pool)
}

func (handler *Handler) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	pool := handler.usecase.GetTransactionPool()
	handler.writeJSON(w, http.StatusOK, pool)
}

func (handler *Handler) GetMineTransactions(w http.ResponseWriter, r *http.Request) {
	if err := handler.usecase.MineTransactions(); err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	handler.GetBlocks(w, r)
}

func (handler *Handler) GetWalletInfo(w http.ResponseWriter, r *http.Request) {
	walletInfo, err := handler.usecase.GetWalletInfo()
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
	}
	handler.writeJSON(w, http.StatusOK, walletInfo)
}
