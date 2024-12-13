package handler

import (
	"context"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/usecase"
)

type Handler struct {
	ctx         context.Context
	configs     *configs.Config
	usecase 	usecase.UseCaseInterface
}

type HandlerInterface interface {
	GetBlocks(w http.ResponseWriter, r *http.Request)
	Mine(w http.ResponseWriter, r *http.Request)
}

var _ HandlerInterface = (*Handler)(nil)

func NewHandler(ctx context.Context, usecase *usecase.UseCase, configs *configs.Config) *Handler {
	return &Handler{
		ctx:        ctx,
		usecase: 	usecase,
		configs:    configs,
	}
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusOK, handler.usecase.GetBlock())
}

func (handler *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	var requestPayload model.Payload

	handler.readJSON(w, r, &requestPayload)

	err := handler.usecase.Mine(requestPayload.Data);
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	handler.GetBlocks(w, r)
}


