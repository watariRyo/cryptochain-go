package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/handler"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
)

type Server struct {
	ctx     context.Context
	handler *handler.Handler
}

func Run() {
	realTimeProvider := &time.RealTimeProvider{}
	ctx := context.Background()
	blockChain := block.NewBlockChain(ctx, realTimeProvider)

	config, err := configs.Load()
	if err != nil {
		log.Panic(err)
	}

	redisClient, err := redis.NewRedisClient(&config.Redis, ctx, blockChain)
	if err != nil {
		log.Panic(err)
	}

	// dependencies
	handler := handler.NewHandler(ctx, blockChain, redisClient, config)

	server := Server{
		ctx:     ctx,
		handler: handler,
	}

	// init broadcast
	broadcastChain, err := json.Marshal(blockChain.GetBlock())
	if err != nil {
		log.Panic(err)
	}

	go redisClient.Publish(ctx, string(redis.BLOCKCHAIN), string(broadcastChain))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: server.Routes(),
	}

	go func() {
		if config.Server.DefaultPort != config.Server.Port {
			err = handler.SyncChain()
			if err != nil {
				log.Panic(err)
			}
		}

		log.Println("Starting service on port", config.Server.Port)

		err = srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}
	}()

	select {}
}
