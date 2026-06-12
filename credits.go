package miget

import (
	"context"
	"net/http"
)

// CreditsService handles the credits endpoints (OpenAPI tag
// "Users - Credits").
type CreditsService struct {
	client *Client
}

// Get returns the authenticated user's credit balance and history.
// Operation: getUserCredits (GET /api/v1/users/me/credits)
func (s *CreditsService) Get(ctx context.Context) (*Credit, error) {
	out := new(Credit)
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/users/me/credits", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
