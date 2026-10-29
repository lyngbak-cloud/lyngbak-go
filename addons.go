package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AddonsService handles the App Add-ons endpoints (OpenAPI tag
// "App Add-ons").
type AddonsService struct {
	client *Client
}

// List returns all add-ons of an application.
// Operation: listAppAddons (GET /api/v1/apps/{uuid}/addons)
func (s *AddonsService) List(ctx context.Context, appUUID string) ([]Addon, error) {
	var out []Addon
	path := fmt.Sprintf("/api/v1/apps/%s/addons", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates an add-on for an application.
// Operation: createAppAddon (POST /api/v1/apps/{uuid}/addons)
func (s *AddonsService) Create(ctx context.Context, appUUID string, req *CreateAppAddonRequest) (*Addon, error) {
	out := new(Addon)
	path := fmt.Sprintf("/api/v1/apps/%s/addons", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single add-on.
// Operation: getAppAddon (GET /api/v1/apps/{uuid}/addons/{id})
func (s *AddonsService) Get(ctx context.Context, appUUID, addonID string) (*Addon, error) {
	out := new(Addon)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates an add-on.
// Operation: updateAppAddon (PUT /api/v1/apps/{uuid}/addons/{id})
func (s *AddonsService) Update(ctx context.Context, appUUID, addonID string, req *UpdateAppAddonRequest) (*Addon, error) {
	out := new(Addon)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes an add-on.
// Operation: deleteAppAddon (DELETE /api/v1/apps/{uuid}/addons/{id})
func (s *AddonsService) Delete(ctx context.Context, appUUID, addonID string) error {
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s", url.PathEscape(appUUID), url.PathEscape(addonID))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdateState starts, stops or restarts an add-on. state is one of the
// ProcessState* constants.
// Operation: updateAppAddonState (PATCH /api/v1/apps/{uuid}/addons/{id}/state)
func (s *AddonsService) UpdateState(ctx context.Context, appUUID, addonID, state string) (*Message, error) {
	body := struct {
		State string `json:"state"`
	}{State: state}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/state", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPatch, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RotatePassword rotates a database add-on's password.
// Operation: rotateAddonPassword (POST /api/v1/apps/{uuid}/addons/{id}/rotate_password)
func (s *AddonsService) RotatePassword(ctx context.Context, appUUID, addonID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/rotate_password", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Backups returns the backup configuration and history of a database add-on.
// Operation: getAppAddonBackups (GET /api/v1/apps/{uuid}/addons/{id}/backups)
func (s *AddonsService) Backups(ctx context.Context, appUUID, addonID string) (*Backup, error) {
	out := new(Backup)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/backups", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RestoreBackup restores a database add-on from a backup. req may be nil to
// restore from the latest backup.
// Operation: restoreAppAddonBackup (POST /api/v1/apps/{uuid}/addons/{id}/restore_backup)
func (s *AddonsService) RestoreBackup(ctx context.Context, appUUID, addonID string, req *RestoreBackupRequest) (*Message, error) {
	if req == nil {
		req = &RestoreBackupRequest{}
	}
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/restore_backup", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ResetDatabase resets a database add-on, destroying all data.
// Operation: resetAppAddonDatabase (POST /api/v1/apps/{uuid}/addons/{id}/reset_database)
func (s *AddonsService) ResetDatabase(ctx context.Context, appUUID, addonID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/reset_database", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateReplica creates a read replica of a database add-on. req may be nil
// to inherit the primary's resources.
// Operation: createAppAddonReplica (POST /api/v1/apps/{uuid}/addons/{id}/create_replica)
func (s *AddonsService) CreateReplica(ctx context.Context, appUUID, addonID string, req *CreateReplicaRequest) (*Addon, error) {
	if req == nil {
		req = &CreateReplicaRequest{}
	}
	out := new(Addon)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/create_replica", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoteReplica promotes a read replica to primary.
// Operation: promoteAppAddonReplica (POST /api/v1/apps/{uuid}/addons/{id}/promote_replica)
func (s *AddonsService) PromoteReplica(ctx context.Context, appUUID, addonID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/promote_replica", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoteExternalReplica promotes an external replica add-on to a standalone
// primary.
// Operation: promoteAppAddonExternalReplica (POST /api/v1/apps/{uuid}/addons/{id}/promote_external)
func (s *AddonsService) PromoteExternalReplica(ctx context.Context, appUUID, addonID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/addons/%s/promote_external", url.PathEscape(appUUID), url.PathEscape(addonID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
