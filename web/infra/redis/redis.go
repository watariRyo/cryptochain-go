package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
)

type CHANNELS string

const (
	TEST       CHANNELS = "TEST"
	BLOCKCHAIN CHANNELS = "BLOCKCHAIN"
)

var channels = []string{string(TEST), string(BLOCKCHAIN)}

// パブリッシャーとサブスクライバーの両方を宣言する理由は、
// PubSubのインスタンスがアプリケーションで両方の役割を果たせるようにするため
type RedisClient struct {
	ctx        context.Context
	publisher  *redis.Client
	subscriber *redis.PubSub
	blockChain *block.BlockChain
}

var _ repository.RedisClientInterface = (*RedisClient)(nil)

func NewRedisClient(cfg *configs.Redis, ctx context.Context, blockChain *block.BlockChain, tm time.TimeProvider) (*RedisClient, error) {
	pub, err := createRedisClient(cfg, ctx)
	if err != nil {
		return nil, err
	}

	redisClient := &RedisClient{
		publisher:  pub,
		blockChain: blockChain,
	}

	redisClient.Subscribe(ctx, tm)

	return redisClient, nil
}

func createRedisClient(cfg *configs.Redis, ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Password,
			DB:       0,
		},
	)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *RedisClient) Subscribe(ctx context.Context, tm time.TimeProvider) {
	// チャネルをサブスクライブ
	c.subscriber = c.publisher.Subscribe(ctx, channels...)

	// メッセージを受信
	go func(pubsub *redis.PubSub) {
		ch := c.subscriber.Channel()
		defer c.subscriber.Close()

		for msg := range ch {
			if BLOCKCHAIN == CHANNELS(msg.Channel) {
				payload := []byte(msg.Payload)
				c.blockChain.UnmarshalAndReplaceBlock(payload, tm)
			}
		}
	}(c.subscriber)
}

func (c *RedisClient) Publish(ctx context.Context, channel, messages string) {
	c.subscriber.Unsubscribe(c.ctx, channel)
	err := c.publisher.Publish(ctx, channel, messages).Err()
	c.subscriber.Subscribe(c.ctx, channel)
	if err != nil {
		logger.Errorf(c.ctx, "Error publishing message: %v\n", err)
		return
	}
}
