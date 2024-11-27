package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/charmbracelet/log"
)

func GenerateRandomID() string {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		log.Error("Error generating player ID: ", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}
