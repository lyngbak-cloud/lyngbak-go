package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// RegistryCredentialsService handles the Container Registry Credentials
// endpoints (OpenAPI tag "Container Registry Credentials").
type RegistryCredentialsService struct {
	client *Client
}

// List returns all container registry credentials in the workspace.
// Operation: listContainerRegistryCredentials (GET /api/v1/container_registry_credentials)
func (s *RegistryCredentialsService) List(ctx context.Context) ([]ContainerRegistryCredential, error) {
	var out []ContainerRegistryCredential
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/container_registry_credentials", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create stores a new container registry credential.
// Operation: createContainerRegistryCredential (POST /api/v1/container_registry_credentials)
func (s *RegistryCredentialsService) Create(ctx context.Context, req *CreateContainerRegistryCredentialRequest) (*ContainerRegistryCredential, error) {
	out := new(ContainerRegistryCredential)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/container_registry_credentials", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single container registry credential.
// Operation: getContainerRegistryCredential (GET /api/v1/container_registry_credentials/{uuid})
func (s *RegistryCredentialsService) Get(ctx context.Context, uuid string) (*ContainerRegistryCredential, error) {
	out := new(ContainerRegistryCredential)
	path := fmt.Sprintf("/api/v1/container_registry_credentials/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a container registry credential.
// Operation: updateContainerRegistryCredential (PUT /api/v1/container_registry_credentials/{uuid})
func (s *RegistryCredentialsService) Update(ctx context.Context, uuid string, req *UpdateContainerRegistryCredentialRequest) (*ContainerRegistryCredential, error) {
	out := new(ContainerRegistryCredential)
	path := fmt.Sprintf("/api/v1/container_registry_credentials/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a container registry credential.
// Operation: deleteContainerRegistryCredential (DELETE /api/v1/container_registry_credentials/{uuid})
func (s *RegistryCredentialsService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/container_registry_credentials/%s", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
