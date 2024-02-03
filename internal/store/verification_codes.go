package store

import (
	"context"
	"fmt"
	"server/internal/types"

	"github.com/redis/go-redis/v9"
)

type VerificationCodes struct {
	redis *redis.Client
}

func NewVerificationCodesStore(redis *redis.Client) *VerificationCodes {
	return &VerificationCodes{
		redis: redis,
	}
}

func (s *VerificationCodes) Create(code, phone string) error {
	err := s.redis.Set(
		context.Background(),
		phone,
		code,
		0,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while storing verification code: %w", err)
	}
	return nil
}

func (s *VerificationCodes) Get(phone string) (string, error) {
	code, err := s.redis.Get(
		context.Background(),
		phone,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrCodeNotFound
		}
		return "", fmt.Errorf("Error while getting verification code: %w", err)
	}
	return code, nil
}

func (s *VerificationCodes) Delete(code, phone string) error {
	err := s.redis.Del(
		context.Background(),
		phone,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while deleting verification code: %w", err)
	}
	return nil
}
