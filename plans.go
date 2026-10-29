package lyngbak

import (
	"context"
	"net/http"
)

// PlansService handles the Plans endpoints (OpenAPI tag "Plans").
type PlansService struct {
	client *Client
}

// List returns all available billing plans.
// Operation: listPlans (GET /api/v1/plans)
func (s *PlansService) List(ctx context.Context) ([]Plan, error) {
	var out []Plan
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/plans", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
