package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CronjobsService handles the App Cron Jobs endpoints (OpenAPI tag
// "App Cron Jobs").
type CronjobsService struct {
	client *Client
}

// List returns all cron jobs of an application.
// Operation: listAppCronjobs (GET /api/v1/apps/{uuid}/cronjobs)
func (s *CronjobsService) List(ctx context.Context, appUUID string) ([]Cronjob, error) {
	var out []Cronjob
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a cron job for an application.
//
// Note: the OpenAPI spec declares the response as the Addon entity, which
// appears to be a spec error; the Cronjob type is decoded here.
// Operation: createAppCronjob (POST /api/v1/apps/{uuid}/cronjobs)
func (s *CronjobsService) Create(ctx context.Context, appUUID string, req *CreateAppCronjobRequest) (*Cronjob, error) {
	out := new(Cronjob)
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single cron job.
// Operation: getAppCronjob (GET /api/v1/apps/{uuid}/cronjobs/{id})
func (s *CronjobsService) Get(ctx context.Context, appUUID, cronjobID string) (*Cronjob, error) {
	out := new(Cronjob)
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs/%s", url.PathEscape(appUUID), url.PathEscape(cronjobID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a cron job.
//
// Note: the OpenAPI spec declares the response as the Addon entity, which
// appears to be a spec error; the Cronjob type is decoded here.
// Operation: updateAppCronjob (PUT /api/v1/apps/{uuid}/cronjobs/{id})
func (s *CronjobsService) Update(ctx context.Context, appUUID, cronjobID string, req *UpdateAppCronjobRequest) (*Cronjob, error) {
	out := new(Cronjob)
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs/%s", url.PathEscape(appUUID), url.PathEscape(cronjobID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a cron job.
// Operation: deleteAppCronjob (DELETE /api/v1/apps/{uuid}/cronjobs/{id})
func (s *CronjobsService) Delete(ctx context.Context, appUUID, cronjobID string) error {
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs/%s", url.PathEscape(appUUID), url.PathEscape(cronjobID))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (s *CronjobsService) StreamLogs(ctx context.Context, appUUID string, id int64, fn func(LogEvent) error) error {
	path := fmt.Sprintf("/api/v1/apps/%s/cronjobs/%d/stream_logs", url.PathEscape(appUUID), id)
	return s.client.streamSSE(ctx, path, fn)
}
