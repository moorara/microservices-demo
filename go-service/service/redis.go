package service

import (
	"time"

	"github.com/go-redis/redis"
)

type (
	// Persister is the interface for a data-backend
	Persister interface {
		Save(key string, data []byte, ttl time.Duration) error
		Load(key string) ([]byte, error)
	}

	// redisPersister implements Persister interface
	redisPersister struct {
		redisClient *redis.Client
	}
)

// NewRedisPersister creates a new Persister for Redis
func NewRedisPersister(redisURL string) Persister {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	return &redisPersister{
		redisClient: redis.NewClient(opts),
	}
}

func (rp *redisPersister) Save(key string, data []byte, ttl time.Duration) error {
	return rp.redisClient.Set(key, data, ttl).Err()
}

func (rp *redisPersister) Load(key string) ([]byte, error) {
	return rp.redisClient.Get(key).Bytes()
}
