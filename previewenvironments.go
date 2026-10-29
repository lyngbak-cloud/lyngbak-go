package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type PreviewEnvironmentsService struct {
	client *Client
}

func (s *PreviewEnvironmentsService) GetConfig(ctx context.Context, appUUID string) (*PreviewEnvironmentConfig, error) {
	out := new(PreviewEnvironmentConfig)
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments/config", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PreviewEnvironmentsService) UpdateConfig(ctx context.Context, appUUID string, req *UpdatePreviewEnvironmentConfigRequest) (*PreviewEnvironmentConfig, error) {
	out := new(PreviewEnvironmentConfig)
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments/config", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PreviewEnvironmentsService) List(ctx context.Context, appUUID string, status string) ([]PreviewEnvironment, error) {
	query := url.Values{}
	if status != "" {
		query.Set("status", status)
	}
	var out []PreviewEnvironment
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PreviewEnvironmentsService) Get(ctx context.Context, appUUID string, id int64) (*PreviewEnvironment, error) {
	out := new(PreviewEnvironment)
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments/%d", url.PathEscape(appUUID), id)
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PreviewEnvironmentsService) Delete(ctx context.Context, appUUID string, id int64) error {
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments/%d", url.PathEscape(appUUID), id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (s *PreviewEnvironmentsService) Redeploy(ctx context.Context, appUUID string, id int64) error {
	path := fmt.Sprintf("/api/v1/apps/%s/preview_environments/%d/redeploy", url.PathEscape(appUUID), id)
	return s.client.do(ctx, http.MethodPost, path, nil, nil, nil)
}
