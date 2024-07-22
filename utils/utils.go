package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateUserNameSuffix() string {
	const LENGTH = 6
	const charset = "0123456789"

	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	random := rand.New(seed)

	result := make([]byte, LENGTH)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}
