package miget

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PortsService handles the App Ports endpoints (OpenAPI tag "App Ports").
type PortsService struct {
	client *Client
}

// List returns all exposed ports of an application.
// Operation: listAppPorts (GET /api/v1/apps/{uuid}/ports)
func (s *PortsService) List(ctx context.Context, appUUID string) ([]Port, error) {
	var out []Port
	path := fmt.Sprintf("/api/v1/apps/%s/ports", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create exposes a new port on an application.
// Operation: createAppPort (POST /api/v1/apps/{uuid}/ports)
func (s *PortsService) Create(ctx context.Context, appUUID string, req *CreateAppPortRequest) (*Port, error) {
	out := new(Port)
	path := fmt.Sprintf("/api/v1/apps/%s/ports", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single port.
// Operation: getAppPort (GET /api/v1/apps/{uuid}/ports/{port_id})
func (s *PortsService) Get(ctx context.Context, appUUID string, portID int64) (*Port, error) {
	out := new(Port)
	path := fmt.Sprintf("/api/v1/apps/%s/ports/%d", url.PathEscape(appUUID), portID)
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes an exposed port from an application.
// Operation: deleteAppPort (DELETE /api/v1/apps/{uuid}/ports/{port_id})
func (s *PortsService) Delete(ctx context.Context, appUUID string, portID int64) error {
	path := fmt.Sprintf("/api/v1/apps/%s/ports/%d", url.PathEscape(appUUID), portID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// ExposePublicly exposes a port to the public internet.
// Operation: exposeAppPortPublicly (PATCH /api/v1/apps/{uuid}/ports/{port_id}/expose_publicly)
func (s *PortsService) ExposePublicly(ctx context.Context, appUUID string, portID int64) (*Port, error) {
	out := new(Port)
	path := fmt.Sprintf("/api/v1/apps/%s/ports/%d/expose_publicly", url.PathEscape(appUUID), portID)
	if err := s.client.do(ctx, http.MethodPatch, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MakePrivate removes the public exposure of a port.
// Operation: makeAppPortPrivate (PATCH /api/v1/apps/{uuid}/ports/{port_id}/make_private)
func (s *PortsService) MakePrivate(ctx context.Context, appUUID string, portID int64) (*Port, error) {
	out := new(Port)
	path := fmt.Sprintf("/api/v1/apps/%s/ports/%d/make_private", url.PathEscape(appUUID), portID)
	if err := s.client.do(ctx, http.MethodPatch, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
