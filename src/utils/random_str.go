package utils

import (
    "crypto/rand"
    "fmt"
)
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(nbytes int) (string, error) {
     // Buffer para 32 bytes aleatorios
    bytes := make([]byte, nbytes)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }

    // Mapear bytes a caracteres del alfabeto
    result := make([]byte, 32)
    for i, b := range bytes {
        result[i] = alphabet[int(b)%len(alphabet)]
    }
    return string(result), nil
}
