package server

import (
	"context"

	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/handler"
)

type Server struct {
	Ctx          context.Context
	TimeProvider *time.RealTimeProvider
	Handler      *handler.Handler
}
