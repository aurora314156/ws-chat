package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"hashed_password,omitempty"`
	FullName       string    `json:"full_name,omitempty"`
}
