package utils

import (
	"fmt"
	"time"
)

func GenerateRandomDigit() string {
	return fmt.Sprint(time.Now().Nanosecond())[:6]
}
