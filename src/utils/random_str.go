package utils

import (
    "crypto/rand"
    "fmt"
    "io"
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

func GenerateRandomOTP(nbytes int) (string, error) {
	const alphabet string = "ABCDEFG1234HIJKLM5678OPQRSTU90VWXYZ"

	// Ensure nbytes is positive to avoid issues with make and for loop
	if nbytes <= 0 {
		return "", fmt.Errorf("nbytes must be a positive integer")
	}

	bytes := make([]byte, nbytes)
	// Read cryptographically secure random bytes
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	result := make([]byte, nbytes) // The result slice should have the same length as nbytes
	for i, b := range bytes {
		result[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(result), nil
}