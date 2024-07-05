package asyclient

import "github.com/hibiken/asynq"

var client *asynq.Client

func Init(redisConnOpt asynq.RedisConnOpt) {
	client = asynq.NewClient(redisConnOpt)
}

func Client() *asynq.Client {
	return client
}
