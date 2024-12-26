package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (server *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/api/blocks", server.handler.GetBlocks)
	mux.Post("/api/mine", server.handler.Mine)
	mux.Post("/api/transact", server.handler.Transact)
	mux.Get("/api/transaction-pool-map", server.handler.GetTransactionPool)
	mux.Get("/api/mine-transactions", server.handler.GetMineTransactions)
	mux.Get("/api/wallet-info", server.handler.GetWalletInfo)

	return mux
}
