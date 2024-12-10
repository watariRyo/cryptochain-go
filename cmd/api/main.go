package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/block"
	"github.com/watariRyo/cryptochain-go/web/handler"
	"github.com/watariRyo/cryptochain-go/web/redis"
	"github.com/watariRyo/cryptochain-go/web/server"
)

func main() {
	realTimeProvider := &tm.RealTimeProvider{}
	ctx := context.Background()
	blockChain := block.NewBlockChain(ctx, realTimeProvider)

	config, err := configs.Load()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewRedisClient(&config.Redis, ctx, blockChain)
	if err != nil {
		panic(err)
	}

	// dependencies
	handler := &handler.Handler{
		BlockChain:  blockChain,
		RedisClient: redisClient,
	}

	server := server.Server{
		Ctx:     ctx,
		Handler: handler,
	}

	log.Println("Starting mail service on port", config.Server.Port)

	// init broadcast
	broadcastChain, err := json.Marshal(blockChain.Block)
	if err != nil {
		panic(err)
	}

	go redisClient.Publish(ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: server.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
