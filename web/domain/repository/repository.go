package repository

type AllRepository struct {
	RedisClient RedisClientInterface
	BlockChain  BlockChainInterface
}

func NewRepository(redisClient RedisClientInterface, blockChain BlockChainInterface) *AllRepository {
	return &AllRepository{
		RedisClient: redisClient,
		BlockChain:  blockChain,
	}
}
