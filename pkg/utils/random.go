package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString generates a random hex string of length n*2
func GenerateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
