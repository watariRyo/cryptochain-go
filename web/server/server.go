package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	"github.com/watariRyo/cryptochain-go/web/handler"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
	"github.com/watariRyo/cryptochain-go/web/infra/redis"
	"github.com/watariRyo/cryptochain-go/web/infra/wallets"
	"github.com/watariRyo/cryptochain-go/web/usecase"
)

type Server struct {
	ctx     context.Context
	handler *handler.Handler
}

func Run() {
	realTimeProvider := &time.RealTimeProvider{}
	ctx := context.Background()
	blockChain := block.NewBlockChain(ctx, realTimeProvider)
	wallet, err := wallets.NewWallet()
	if err != nil {
		log.Panic(err)
	}
	wallets := wallets.NewWallets(wallet, nil)

	configs, err := configs.Load()
	if err != nil {
		log.Panic(err)
	}

	redisClient, err := redis.NewRedisClient(&configs.Redis, ctx, blockChain, wallets, realTimeProvider)
	if err != nil {
		log.Panic(err)
	}

	// dependencies
	repo := repository.NewRepository(redisClient, blockChain, wallets)

	usecase := usecase.NewUseCase(ctx, realTimeProvider, repo, configs)

	handler := handler.NewHandler(ctx, usecase, configs)

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
		Addr:    fmt.Sprintf(":%s", configs.Server.Port),
		Handler: server.Routes(),
	}

	go func() {
		if configs.Server.DefaultPort != configs.Server.Port {
			err = usecase.SyncWithRootState()
			if err != nil {
				log.Panic(err)
			}
		}

		log.Println("Starting service on port", configs.Server.Port)

		err = srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}
	}()

	select {}
}
