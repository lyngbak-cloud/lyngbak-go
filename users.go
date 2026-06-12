package miget

import (
	"context"
	"net/http"
)

// UsersService handles the Users endpoints (OpenAPI tag "Users").
type UsersService struct {
	client *Client
}

// Me returns the authenticated user.
// Operation: getCurrentUser (GET /api/v1/users/me)
func (s *UsersService) Me(ctx context.Context) (*User, error) {
	out := new(User)
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/users/me", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
