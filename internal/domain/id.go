package domain

import "github.com/google/uuid"

// GenerateID generates a unique ID that can be used as an identifier for a domain.
func GenerateID() string {
	return uuid.New().String()
}
