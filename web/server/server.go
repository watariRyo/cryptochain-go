package server

import (
	"context"

	"github.com/watariRyo/cryptochain-go/web/handler"
)

type Server struct {
	Ctx     context.Context
	Handler *handler.Handler
}
