package tool

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ToString(val any, fallback string) string {
	if val == nil {
		return fallback
	}
	s, ok := val.(string)
	if !ok {
		return fallback
	}
	return s
}

func ConvertUTCToISO(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}

func GenUUID() uuid.UUID {
	return uuid.New()
}

func HashedPassword(password []byte) (hashedPassword []byte, err error) {
	hashedPassword, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return hashedPassword, err
}
