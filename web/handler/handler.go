package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

type Handler struct {
	ctx         context.Context
	blockChain  block.BlockChainInterface
	redisClient redis.RedisClientInterface
	configs     *configs.Config
}

type HandlerInterface interface {
	GetBlocks(w http.ResponseWriter, r *http.Request)
	Mine(w http.ResponseWriter, r *http.Request)
	SyncChain() error
}

var _ HandlerInterface = (*Handler)(nil)

func NewHandler(ctx context.Context, blockChain *block.BlockChain, redisClient *redis.RedisClient, configs *configs.Config) *Handler {
	return &Handler{
		ctx:         ctx,
		blockChain:  blockChain,
		redisClient: redisClient,
		configs:     configs,
	}
}

func (handler *Handler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	handler.writeJSON(w, http.StatusOK, handler.blockChain.GetBlock())
}

func (handler *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	var requestPayload payload

	handler.readJSON(w, r, &requestPayload)

	handler.blockChain.AddBlock(requestPayload.Data)

	broadcastChain, err := json.Marshal(handler.blockChain.GetBlock())
	if err != nil {
		handler.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	go handler.redisClient.Publish(handler.ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	handler.GetBlocks(w, r)
}

func (handler *Handler) SyncChain() error {
	// Route Node
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blocks", handler.configs.Host), nil)
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

	handler.blockChain.UnmarshalAndReplaceBlock(payload)

	return nil
}
