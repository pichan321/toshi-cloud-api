package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))

	return  hex.EncodeToString(h.Sum(nil))
}

func GenerateUuid() string {
	return uuid.New().String()
}

func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}