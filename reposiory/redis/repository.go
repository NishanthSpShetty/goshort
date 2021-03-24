package redis

import (
	"github.com/go-redis/redis"
)

type redisRepository struct {
	client *redis.Client
}
