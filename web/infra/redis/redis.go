package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/watariRyo/cryptochain-go/configs"
	"github.com/watariRyo/cryptochain-go/internal/logger"
	"github.com/watariRyo/cryptochain-go/internal/time"
	"github.com/watariRyo/cryptochain-go/web/domain/model"
	"github.com/watariRyo/cryptochain-go/web/domain/repository"
	"github.com/watariRyo/cryptochain-go/web/infra/block"
	"github.com/watariRyo/cryptochain-go/web/infra/wallets"
)

type CHANNELS string

const (
	TEST        CHANNELS = "TEST"
	BLOCKCHAIN  CHANNELS = "BLOCKCHAIN"
	TRANSACTION CHANNELS = "TRANSACTION"
)

var channels = []string{string(TEST), string(BLOCKCHAIN), string(TRANSACTION)}

// パブリッシャーとサブスクライバーの両方を宣言する理由は、
// PubSubのインスタンスがアプリケーションで両方の役割を果たせるようにするため
type RedisClient struct {
	publisher  *redis.Client
	subscriber *redis.PubSub
	blockChain *block.BlockChain
	wallets    *wallets.Wallets
}

var _ repository.RedisClientInterface = (*RedisClient)(nil)

func NewRedisClient(cfg *configs.Redis, ctx context.Context, blockChain *block.BlockChain, wallets *wallets.Wallets, tm time.TimeProvider) (*RedisClient, error) {
	pub, err := createRedisClient(cfg, ctx)
	if err != nil {
		return nil, err
	}

	redisClient := &RedisClient{
		publisher:  pub,
		blockChain: blockChain,
		wallets:    wallets,
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
				validTransactionDataFn := c.wallets.ValidTransactionData
				c.blockChain.UnmarshalAndReplaceBlock(ctx, payload, tm, c.wallets.ClearBlockChainTransactions, validTransactionDataFn)
			}
			if TRANSACTION == CHANNELS(msg.Channel) {
				payload := []byte(msg.Payload)
				var payloadTransaction *model.Transaction
				if err := json.Unmarshal(payload, &payloadTransaction); err != nil {
					logger.Errorf(ctx, "Could not unmarshal block chain. %v", err)
				}

				c.wallets.SetTransaction(payloadTransaction)
			}
		}
	}(c.subscriber)
}

func (c *RedisClient) Publish(ctx context.Context, channel, messages string) {
	c.subscriber.Unsubscribe(ctx, channel)
	err := c.publisher.Publish(ctx, channel, messages).Err()
	c.subscriber.Subscribe(ctx, channel)
	if err != nil {
		logger.Errorf(ctx, "Error publishing message: %v\n", err)
		return
	}
}
