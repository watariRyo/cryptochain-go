package repository

import (
	"context"

	"github.com/watariRyo/cryptochain-go/internal/time"
)

type RedisClientInterface interface {
	Subscribe(ctx context.Context, tm time.TimeProvider)
	Publish(ctx context.Context, channel, messages string)
}
