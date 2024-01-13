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

func (oc *OTPCodes) Set(codeKey string) error {
	err := oc.redis.SetEx(
		context.Background(),
		codeKey,
		"",
		time.Minute*5,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while setting OTP code: %w", err)
	}
	return nil
}

func (oc *OTPCodes) Get(codeKey string) (string, error) {
	value, err := oc.redis.Get(
		context.Background(),
		codeKey,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", types.ErrCodeNotFound
		}
		return "", fmt.Errorf("Error while getting OPT: %w", err)
	}
	return value, nil
}

func (oc *OTPCodes) SetVerificationProof(verifiedNumber, proofId string) error {
	err := oc.redis.SetEx(
		context.Background(),
		proofId,
		verifiedNumber,
		time.Minute*5,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while setting verification prooof: %w", err)
	}
	return nil
}

func (oc *OTPCodes) DeleteVerificationProof(proofId string) error {
	err := oc.redis.Del(
		context.Background(),
		proofId,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while deleting verification proof: %w", err)
	}
	return nil
}

func (oc *OTPCodes) Delete(key string) error {
	err := oc.redis.Del(
		context.Background(),
		key,
	).Err()
	if err != nil {
		return fmt.Errorf("Error while deleting OTP code: %w", err)
	}
	return nil
}
