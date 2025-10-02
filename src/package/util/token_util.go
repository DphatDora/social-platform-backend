package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Generates a random token of specified lengt
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
