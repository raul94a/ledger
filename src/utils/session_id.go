package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"

)

func CreateSessionId() (string, error) {
	uuid := uuid.NewString()
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(uuid)); err != nil{
		return "",err
	}
	hashBytes := hasher.Sum(nil)
	sessionId := hex.EncodeToString(hashBytes)
	return sessionId,nil
}