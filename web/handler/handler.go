package handler

import (
	"net/http"

	"github.com/watariRyo/cryptochain-go/web/block"
)

type Handler struct {
	BlockChain *block.BlockChain
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusAccepted, handler.BlockChain.Block)
}
