package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ResourcesService handles the Resources endpoints (OpenAPI tag
// "Resources"). A resource (also called a Miget) is a unit of compute that
// apps, services and buckets are provisioned on.
type ResourcesService struct {
	client *Client
}

// List returns all resources in the workspace.
// Operation: listResources (GET /api/v1/resources)
func (s *ResourcesService) List(ctx context.Context) ([]Resource, error) {
	var out []Resource
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/resources", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create provisions a new resource.
// Operation: createResource (POST /api/v1/resources)
func (s *ResourcesService) Create(ctx context.Context, req *CreateResourceRequest) (*Resource, error) {
	out := new(Resource)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/resources", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single resource by UUID.
// Operation: getResource (GET /api/v1/resources/{uuid})
func (s *ResourcesService) Get(ctx context.Context, uuid string) (*Resource, error) {
	out := new(Resource)
	path := fmt.Sprintf("/api/v1/resources/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update changes a resource's plan or components.
// Operation: updateResource (PUT /api/v1/resources/{uuid})
func (s *ResourcesService) Update(ctx context.Context, uuid string, req *UpdateResourceRequest) (*Resource, error) {
	out := new(Resource)
	path := fmt.Sprintf("/api/v1/resources/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a resource.
// Operation: deleteResource (DELETE /api/v1/resources/{uuid})
func (s *ResourcesService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/resources/%s", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdateLabels replaces a resource's labels.
// Operation: updateResourceLabels (PATCH /api/v1/resources/{uuid}/labels)
func (s *ResourcesService) UpdateLabels(ctx context.Context, uuid string, labels []string) (*Resource, error) {
	body := struct {
		Labels []string `json:"labels"`
	}{Labels: labels}
	out := new(Resource)
	path := fmt.Sprintf("/api/v1/resources/%s/labels", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPatch, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}
