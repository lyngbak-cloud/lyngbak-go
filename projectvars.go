package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ProjectVarsService handles the Project Environment Variables endpoints
// (OpenAPI tag "Project Environment Variables").
type ProjectVarsService struct {
	client *Client
}

// List returns all environment variables of a project.
// Operation: listProjectVars (GET /api/v1/projects/{project_id}/vars)
func (s *ProjectVarsService) List(ctx context.Context, projectID string) ([]Variable, error) {
	var out []Variable
	path := fmt.Sprintf("/api/v1/projects/%s/vars", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create sets a new environment variable on a project.
// Operation: createProjectVar (POST /api/v1/projects/{project_id}/vars)
func (s *ProjectVarsService) Create(ctx context.Context, projectID, key, value string) error {
	body := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{Key: key, Value: value}
	path := fmt.Sprintf("/api/v1/projects/%s/vars", url.PathEscape(projectID))
	return s.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Update updates an environment variable on a project.
// Operation: updateProjectVar (PUT /api/v1/projects/{project_id}/vars)
func (s *ProjectVarsService) Update(ctx context.Context, projectID, key, value string) error {
	body := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{Key: key, Value: value}
	path := fmt.Sprintf("/api/v1/projects/%s/vars", url.PathEscape(projectID))
	return s.client.do(ctx, http.MethodPut, path, nil, body, nil)
}

// Delete removes an environment variable from a project by key.
// Operation: deleteProjectVar (DELETE /api/v1/projects/{project_id}/vars)
func (s *ProjectVarsService) Delete(ctx context.Context, projectID, key string) error {
	query := url.Values{}
	query.Set("key", key)
	path := fmt.Sprintf("/api/v1/projects/%s/vars", url.PathEscape(projectID))
	return s.client.do(ctx, http.MethodDelete, path, query, nil, nil)
}
