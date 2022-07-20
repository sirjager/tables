package tokens

import "time"

// Builder is an interface for managing tokens
type Builder interface {
	// Create Token if token for specific id and duration
	CreateToken(user string, duration time.Duration) (string, *Payload, error)

	// Verify Token if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
