package utils

import (
	"fmt"
	"time"
)

func GenerateVerificationCode() string {
  return fmt.Sprint(time.Now().Nanosecond())[:6]
}
