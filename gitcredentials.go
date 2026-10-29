package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GitCredentialsService struct {
	client *Client
}

func (s *GitCredentialsService) List(ctx context.Context) ([]GitCredential, error) {
	var out []GitCredential
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/git_credentials", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *GitCredentialsService) Get(ctx context.Context, uuid string) (*GitCredential, error) {
	out := new(GitCredential)
	path := fmt.Sprintf("/api/v1/git_credentials/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
