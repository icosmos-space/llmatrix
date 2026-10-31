package repository

import "context"

// User is a proxy user account.
type User struct {
	ID    string
	Email string
}

// UserRepository persists users.
type UserRepository interface {
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}
