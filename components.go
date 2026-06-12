package miget

import (
	"context"
	"net/http"
)

// ComponentsService handles the Components endpoints (OpenAPI tag
// "Components").
type ComponentsService struct {
	client *Client
}

// List returns all available resource components (extra RAM, CPU, disk).
// Operation: listComponents (GET /api/v1/components)
func (s *ComponentsService) List(ctx context.Context) ([]Component, error) {
	var out []Component
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/components", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
