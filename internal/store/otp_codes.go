package store

import (
	"context"
	"fmt"
	"server/internal/types"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPCodes struct {
	redis *redis.Client
}

func NewOTPCodes(redis *redis.Client) *OTPCodes {
	return &OTPCodes{
		redis: redis,
	}
}

func (oc *OTPCodes) Set(code, sent_to string) error {
	err := oc.redis.SetEx(
		context.Background(),
		fmt.Sprintf("%s:%s", code, sent_to),
		"",
		time.Minute*5,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while setting OTP code: %w", err)
	}
	return nil
}

func (oc *OTPCodes) Get(code, sent_to string) (string, error) {
	value, err := oc.redis.Get(
		context.Background(),
		fmt.Sprintf("%s:%s", code, sent_to),
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrCodeNotFound
		}
		return "", fmt.Errorf("Error while getting OPT: %w", err)
	}
	return value, nil
}

func (oc *OTPCodes) Delete(code, sent_to string) error {
	err := oc.redis.Del(
		context.Background(),
		fmt.Sprintf("%s:%s", code, sent_to),
	).Err()
	if err != nil {
		return fmt.Errorf("Error while deleting OTP code: %w", err)
	}
	return nil
}
