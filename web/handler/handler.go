package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/web/block"
	"github.com/watariRyo/cryptochain-go/web/redis"
)

type Handler struct {
	BlockChain  block.BlockChainInterface
	RedisClient redis.RedisClientInterface
	Configs     *configs.Config
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusOK, handler.BlockChain.GetBlock())
}

func (handler *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	var requestPayload payload

	handler.readJSON(w, r, &requestPayload)

	handler.BlockChain.AddBlock(requestPayload.Data)

	broadcastChain, err := json.Marshal(handler.BlockChain.GetBlock())
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	go handler.RedisClient.Publish(context.TODO(), string(redis.BLOCKCHAIN), string(broadcastChain))

	handler.GetBlocks(w, r)
}

func (handler *Handler) SyncChain() error {
	// Route Node
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blocks", handler.Configs.Host), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Request Error. StatusCode: %d", response.StatusCode)
	}

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	handler.BlockChain.UnmarshalAndReplaceBlock(payload)

	return nil
}
