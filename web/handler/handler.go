package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/watariRyo/cryptochain-go/web/block"
	"github.com/watariRyo/cryptochain-go/web/redis"
)

type Handler struct {
	BlockChain  *block.BlockChain
	RedisClient *redis.RedisClient
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusAccepted, handler.BlockChain.Block)
}

func (handler *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	var requestPayload payload

	handler.readJSON(w, r, &requestPayload)

	handler.BlockChain.AddBlock(requestPayload.Data)

	broadcastChain, err := json.Marshal(handler.BlockChain.Block)
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	go handler.RedisClient.Publish(context.TODO(), string(redis.BLOCKCHAIN), string(broadcastChain))

	handler.GetBlocks(w, r)
}
