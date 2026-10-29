package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ServicesService handles the Services endpoints (OpenAPI tag "Services").
// Services are standalone databases or shared storage that can be mounted
// into multiple apps.
type ServicesService struct {
	client *Client
}

// List returns all services in the workspace.
// Operation: listServices (GET /api/v1/services)
func (s *ServicesService) List(ctx context.Context) ([]Service, error) {
	var out []Service
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/services", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new service.
// Operation: createService (POST /api/v1/services)
func (s *ServicesService) Create(ctx context.Context, req *CreateServiceRequest) (*Service, error) {
	out := new(Service)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/services", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single service.
// Operation: getService (GET /api/v1/services/{id})
func (s *ServicesService) Get(ctx context.Context, serviceID string) (*Service, error) {
	out := new(Service)
	path := fmt.Sprintf("/api/v1/services/%s", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a service.
// Operation: updateService (PUT /api/v1/services/{id})
func (s *ServicesService) Update(ctx context.Context, serviceID string, req *UpdateServiceRequest) (*Service, error) {
	out := new(Service)
	path := fmt.Sprintf("/api/v1/services/%s", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a service.
// Operation: deleteService (DELETE /api/v1/services/{id})
func (s *ServicesService) Delete(ctx context.Context, serviceID string) error {
	path := fmt.Sprintf("/api/v1/services/%s", url.PathEscape(serviceID))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdateState starts, stops or restarts a service. state is one of the
// ProcessState* constants.
// Operation: updateServiceState (PATCH /api/v1/services/{id}/state)
func (s *ServicesService) UpdateState(ctx context.Context, serviceID, state string) (*Message, error) {
	body := struct {
		State string `json:"state"`
	}{State: state}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/state", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPatch, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MountApp mounts an application to the service.
// Operation: mountAppToService (POST /api/v1/services/{id}/mount_app)
func (s *ServicesService) MountApp(ctx context.Context, serviceID string, req *MountAppRequest) (*Service, error) {
	out := new(Service)
	path := fmt.Sprintf("/api/v1/services/%s/mount_app", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnmountApp unmounts an application from the service.
// Operation: unmountAppFromService (POST /api/v1/services/{id}/unmount_app)
func (s *ServicesService) UnmountApp(ctx context.Context, serviceID, appID string) (*Service, error) {
	body := struct {
		AppID string `json:"app_id"`
	}{AppID: appID}
	out := new(Service)
	path := fmt.Sprintf("/api/v1/services/%s/unmount_app", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RotatePassword rotates a database service's password.
// Operation: rotateServicePassword (POST /api/v1/services/{id}/rotate_password)
func (s *ServicesService) RotatePassword(ctx context.Context, serviceID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/rotate_password", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Backups returns the backup configuration and history of a database
// service.
// Operation: getServiceBackups (GET /api/v1/services/{id}/backups)
func (s *ServicesService) Backups(ctx context.Context, serviceID string) (*Backup, error) {
	out := new(Backup)
	path := fmt.Sprintf("/api/v1/services/%s/backups", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RestoreBackup restores a database service from a backup. req may be nil
// to restore from the latest backup.
// Operation: restoreServiceBackup (POST /api/v1/services/{id}/restore_backup)
func (s *ServicesService) RestoreBackup(ctx context.Context, serviceID string, req *RestoreBackupRequest) (*Message, error) {
	if req == nil {
		req = &RestoreBackupRequest{}
	}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/restore_backup", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ResetDatabase resets a database service, destroying all data.
// Operation: resetServiceDatabase (POST /api/v1/services/{id}/reset_database)
func (s *ServicesService) ResetDatabase(ctx context.Context, serviceID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/reset_database", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateReplica creates a read replica of a database service. req may be
// nil to inherit the primary's resources.
// Operation: createServiceReplica (POST /api/v1/services/{id}/create_replica)
func (s *ServicesService) CreateReplica(ctx context.Context, serviceID string, req *CreateReplicaRequest) (*Service, error) {
	if req == nil {
		req = &CreateReplicaRequest{}
	}
	out := new(Service)
	path := fmt.Sprintf("/api/v1/services/%s/create_replica", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoteReplica promotes a read replica to primary.
// Operation: promoteServiceReplica (POST /api/v1/services/{id}/promote_replica)
func (s *ServicesService) PromoteReplica(ctx context.Context, serviceID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/promote_replica", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoteExternalReplica promotes an external replica service to a
// standalone primary.
// Operation: promoteServiceExternalReplica (POST /api/v1/services/{id}/promote_external)
func (s *ServicesService) PromoteExternalReplica(ctx context.Context, serviceID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/services/%s/promote_external", url.PathEscape(serviceID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
