package miget

import (
	"context"
	"net/http"
)

// RegionsService handles the Regions endpoints (OpenAPI tag "Regions").
type RegionsService struct {
	client *Client
}

// List returns all available deployment regions.
// Operation: listRegions (GET /api/v1/regions)
func (s *RegionsService) List(ctx context.Context) ([]Region, error) {
	var out []Region
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/regions", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
