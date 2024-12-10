package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	tm "github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/block"
	"github.com/watariRyo/cryptochain-go/web/handler"
	"github.com/watariRyo/cryptochain-go/web/server"
)

const webPort = "8080"

func main() {
	realTimeProvider := &tm.RealTimeProvider{}
	ctx := context.Background()
	blockChain := block.NewBlockChain(ctx, realTimeProvider)

	// dependencies
	handler := &handler.Handler{
		BlockChain: blockChain,
	}

	server := server.Server{
		Ctx:          ctx,
		TimeProvider: realTimeProvider,
		Handler:      handler,
	}

	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: server.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
