package miget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AppsService handles the Apps endpoints (OpenAPI tag "Apps").
type AppsService struct {
	client *Client
}

// List returns all applications in the workspace.
// Operation: listApps (GET /api/v1/apps)
func (s *AppsService) List(ctx context.Context) ([]App, error) {
	var out []App
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/apps", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new application.
// Operation: createApp (POST /api/v1/apps)
func (s *AppsService) Create(ctx context.Context, req *CreateAppRequest) (*App, error) {
	out := new(App)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/apps", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single application by UUID.
// Operation: getApp (GET /api/v1/apps/{uuid})
func (s *AppsService) Get(ctx context.Context, uuid string) (*App, error) {
	out := new(App)
	path := fmt.Sprintf("/api/v1/apps/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates an application.
// Operation: updateApp (PUT /api/v1/apps/{uuid})
func (s *AppsService) Update(ctx context.Context, uuid string, req *UpdateAppRequest) (*App, error) {
	out := new(App)
	path := fmt.Sprintf("/api/v1/apps/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes an application.
// Operation: deleteApp (DELETE /api/v1/apps/{uuid})
func (s *AppsService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/apps/%s", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdateState starts, stops or restarts an application. state is one of the
// AppStateSchedule* constants.
// Operation: updateAppState (PATCH /api/v1/apps/{uuid}/state)
func (s *AppsService) UpdateState(ctx context.Context, uuid, state string) (*Message, error) {
	body := struct {
		State string `json:"state"`
	}{State: state}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/state", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPatch, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateSecurity updates an application's security settings.
// Operation: updateAppSecurity (PUT /api/v1/apps/{uuid}/security)
func (s *AppsService) UpdateSecurity(ctx context.Context, uuid string, req *UpdateAppSecurityRequest) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/security", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Clone clones an application.
// Operation: cloneApp (POST /api/v1/apps/{uuid}/clone)
func (s *AppsService) Clone(ctx context.Context, uuid string, req *CloneAppRequest) (*App, error) {
	out := new(App)
	path := fmt.Sprintf("/api/v1/apps/%s/clone", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Deploy triggers a deployment of the application. req may be nil to deploy
// with default settings.
// Operation: deployApp (POST /api/v1/apps/{uuid}/deploy)
func (s *AppsService) Deploy(ctx context.Context, uuid string, req *DeployAppRequest) (*Message, error) {
	if req == nil {
		req = &DeployAppRequest{}
	}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/deploy", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateDeployment changes an application's deployment method and
// configuration.
// Operation: updateAppDeployment (PUT /api/v1/apps/{uuid}/deployment)
func (s *AppsService) UpdateDeployment(ctx context.Context, uuid string, req *UpdateAppDeploymentRequest) (*App, error) {
	out := new(App)
	path := fmt.Sprintf("/api/v1/apps/%s/deployment", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Activity returns the application's activity log. opts may be nil.
//
// Note: the OpenAPI spec declares the response as a single Activity entity,
// but the endpoint is paginated (page/limit), so a list is decoded here.
// Flagged for validation against the live API.
// Operation: getAppActivity (GET /api/v1/apps/{uuid}/activity)
func (s *AppsService) Activity(ctx context.Context, uuid string, opts *ActivityListOptions) ([]Activity, error) {
	query := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
	}
	var out []Activity
	path := fmt.Sprintf("/api/v1/apps/%s/activity", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateScalingProfile updates an application's scaling configuration.
// Operation: updateAppScalingProfile (PUT /api/v1/apps/{uuid}/scaling_profile)
func (s *AppsService) UpdateScalingProfile(ctx context.Context, uuid string, req *UpdateAppScalingProfileRequest) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/scaling_profile", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateHealthChecks updates an application's health check configuration.
// Operation: updateAppHealthChecks (PUT /api/v1/apps/{uuid}/health_checks)
func (s *AppsService) UpdateHealthChecks(ctx context.Context, uuid string, req *UpdateAppHealthChecksRequest) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/health_checks", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}
