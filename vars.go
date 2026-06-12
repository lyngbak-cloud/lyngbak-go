package miget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// VarsService handles the App Environment Variables endpoints (OpenAPI tag
// "App Environment Variables").
type VarsService struct {
	client *Client
}

// List returns all environment variables of an application.
// Operation: listAppVars (GET /api/v1/apps/{uuid}/vars)
func (s *VarsService) List(ctx context.Context, appUUID string) ([]Variable, error) {
	var out []Variable
	path := fmt.Sprintf("/api/v1/apps/%s/vars", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create sets a new environment variable on an application.
// Operation: createAppVar (POST /api/v1/apps/{uuid}/vars)
func (s *VarsService) Create(ctx context.Context, appUUID, key, value string) error {
	body := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{Key: key, Value: value}
	path := fmt.Sprintf("/api/v1/apps/%s/vars", url.PathEscape(appUUID))
	return s.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Update updates an environment variable on an application.
// Operation: updateAppVar (PUT /api/v1/apps/{uuid}/vars)
func (s *VarsService) Update(ctx context.Context, appUUID string, req *UpdateAppVarRequest) error {
	path := fmt.Sprintf("/api/v1/apps/%s/vars", url.PathEscape(appUUID))
	return s.client.do(ctx, http.MethodPut, path, nil, req, nil)
}

// Delete removes an environment variable from an application by key.
// Operation: deleteAppVar (DELETE /api/v1/apps/{uuid}/vars)
func (s *VarsService) Delete(ctx context.Context, appUUID, key string) error {
	query := url.Values{}
	query.Set("key", key)
	path := fmt.Sprintf("/api/v1/apps/%s/vars", url.PathEscape(appUUID))
	return s.client.do(ctx, http.MethodDelete, path, query, nil, nil)
}
