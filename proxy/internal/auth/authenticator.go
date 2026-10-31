package auth

import "context"

// Authenticator validates credentials for incoming requests.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*Identity, error)
}

// Identity is an authenticated principal.
type Identity struct {
	UserID string
	KeyID  string
}

// StaticAuthenticator accepts a single configured API key.
type StaticAuthenticator struct {
	APIKey string
	UserID string
}

// Authenticate implements Authenticator.
func (a *StaticAuthenticator) Authenticate(_ context.Context, token string) (*Identity, error) {
	if a.APIKey == "" || token != a.APIKey {
		return nil, ErrUnauthorized
	}
	return &Identity{UserID: a.UserID, KeyID: "static"}, nil
}
