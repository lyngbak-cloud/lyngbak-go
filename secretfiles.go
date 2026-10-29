package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type SecretFilesService struct {
	client *Client
}

func (s *SecretFilesService) List(ctx context.Context, appUUID string) ([]SecretFile, error) {
	var out []SecretFile
	path := fmt.Sprintf("/api/v1/apps/%s/secret_files", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SecretFilesService) Create(ctx context.Context, appUUID string, req *SecretFileRequest) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/apps/%s/secret_files", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SecretFilesService) Get(ctx context.Context, appUUID, name string) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/apps/%s/secret_files/%s", url.PathEscape(appUUID), url.PathEscape(name))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SecretFilesService) Update(ctx context.Context, appUUID, name string, req *SecretFileRequest) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/apps/%s/secret_files/%s", url.PathEscape(appUUID), url.PathEscape(name))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SecretFilesService) Delete(ctx context.Context, appUUID, name string) error {
	path := fmt.Sprintf("/api/v1/apps/%s/secret_files/%s", url.PathEscape(appUUID), url.PathEscape(name))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

type ProjectSecretFilesService struct {
	client *Client
}

func (s *ProjectSecretFilesService) List(ctx context.Context, projectID string) ([]SecretFile, error) {
	var out []SecretFile
	path := fmt.Sprintf("/api/v1/projects/%s/secret_files", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProjectSecretFilesService) Create(ctx context.Context, projectID string, req *SecretFileRequest) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/projects/%s/secret_files", url.PathEscape(projectID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProjectSecretFilesService) Get(ctx context.Context, projectID, name string) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/projects/%s/secret_files/%s", url.PathEscape(projectID), url.PathEscape(name))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProjectSecretFilesService) Update(ctx context.Context, projectID, name string, req *SecretFileRequest) (*SecretFile, error) {
	out := new(SecretFile)
	path := fmt.Sprintf("/api/v1/projects/%s/secret_files/%s", url.PathEscape(projectID), url.PathEscape(name))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProjectSecretFilesService) Delete(ctx context.Context, projectID, name string) error {
	path := fmt.Sprintf("/api/v1/projects/%s/secret_files/%s", url.PathEscape(projectID), url.PathEscape(name))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
