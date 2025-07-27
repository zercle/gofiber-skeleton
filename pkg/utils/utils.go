package utils

import (
	"math/rand"
	"time"
)

var (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateShortCode generates a random short code of a given length.
// If length is less than or equal to zero, it returns an empty string.
func GenerateShortCode(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
