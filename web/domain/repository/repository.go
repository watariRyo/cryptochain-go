package repository

type AllRepository struct {
	RedisClient RedisClientInterface
	BlockChain  BlockChainInterface
	Wallets     WalletsInterface
}

func NewRepository(redisClient RedisClientInterface, blockChain BlockChainInterface, wallets WalletsInterface) *AllRepository {
	return &AllRepository{
		RedisClient: redisClient,
		BlockChain:  blockChain,
		Wallets:     wallets,
	}
}
