package store

import "github.com/redis/go-redis/v9"

type Sessions struct {
	redis *redis.Client
}

func NewSessionsStore(redis *redis.Client) *Sessions {
	return &Sessions{
		redis: redis,
	}
}

func (s *Sessions) Create() error {
	return nil
}
