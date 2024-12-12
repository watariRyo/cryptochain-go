package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/block"
	"github.com/watariRyo/cryptochain-go/web/handler"
	"github.com/watariRyo/cryptochain-go/web/redis"
)

type Server struct {
	Ctx     context.Context
	Handler *handler.Handler
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
	handler := &handler.Handler{
		BlockChain:  blockChain,
		RedisClient: redisClient,
		Configs:     config,
	}

	server := Server{
		Ctx:     ctx,
		Handler: handler,
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
