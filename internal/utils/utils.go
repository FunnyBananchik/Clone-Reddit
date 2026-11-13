package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "000000000000000000000000"
	}
	return hex.EncodeToString(b)
}
