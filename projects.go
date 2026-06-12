package miget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ProjectsService handles the Projects endpoints (OpenAPI tag "Projects").
type ProjectsService struct {
	client *Client
}

// List returns all projects in the workspace.
// Operation: listProjects (GET /api/v1/projects)
func (s *ProjectsService) List(ctx context.Context) ([]Project, error) {
	var out []Project
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/projects", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new project.
// Operation: createProject (POST /api/v1/projects)
func (s *ProjectsService) Create(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	out := new(Project)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/projects", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single project.
// Operation: getProject (GET /api/v1/projects/{project_id})
func (s *ProjectsService) Get(ctx context.Context, projectID string) (*Project, error) {
	out := new(Project)
	path := fmt.Sprintf("/api/v1/projects/%s", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a project.
// Operation: updateProject (PUT /api/v1/projects/{project_id})
func (s *ProjectsService) Update(ctx context.Context, projectID string, req *UpdateProjectRequest) (*Project, error) {
	out := new(Project)
	path := fmt.Sprintf("/api/v1/projects/%s", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a project.
// Operation: deleteProject (DELETE /api/v1/projects/{project_id})
func (s *ProjectsService) Delete(ctx context.Context, projectID string) error {
	path := fmt.Sprintf("/api/v1/projects/%s", url.PathEscape(projectID))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ListApps returns all applications in a project.
// Operation: listProjectApps (GET /api/v1/projects/{project_id}/apps)
func (s *ProjectsService) ListApps(ctx context.Context, projectID string) ([]App, error) {
	var out []App
	path := fmt.Sprintf("/api/v1/projects/%s/apps", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
