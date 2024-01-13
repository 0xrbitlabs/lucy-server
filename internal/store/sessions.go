package store

import (
	"context"
	"fmt"
	"server/internal/types"
	"time"

	"github.com/redis/go-redis/v9"
)

type Sessions struct {
	redis *redis.Client
}

func NewSessionsStore(redis *redis.Client) *Sessions {
	return &Sessions{
		redis: redis,
	}
}

func (s *Sessions) Create(sessionId, userId string) error {
	err := s.redis.SetEx(
		context.Background(),
		sessionId,
		userId,
		time.Hour*24*31,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while creating session: %w", err)
	}
	return nil
}

func (s *Sessions) Get(sessionId string) (string, error) {
	sessionUser, err := s.redis.Get(
		context.Background(),
		sessionId,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrSessionNotFound
		}
		return "", fmt.Errorf("Error while getting session: %w", err)
	}
	return sessionUser, nil
}
