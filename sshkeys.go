package lyngbak

import (
	"context"
	"fmt"
	"net/http"
)

// SSHKeysService handles the SSH key endpoints (OpenAPI tag
// "Users - SSH Keys").
type SSHKeysService struct {
	client *Client
}

// List returns the authenticated user's SSH keys.
// Operation: listSshKeys (GET /api/v1/users/me/ssh_keys)
func (s *SSHKeysService) List(ctx context.Context) ([]SSHKey, error) {
	var out []SSHKey
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/users/me/ssh_keys", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create adds an SSH public key in OpenSSH format, e.g.
// "ssh-ed25519 AAAA... user@host".
// Operation: createSshKey (POST /api/v1/users/me/ssh_keys)
func (s *SSHKeysService) Create(ctx context.Context, publicKey string) (*SSHKey, error) {
	body := struct {
		PublicKey string `json:"public_key"`
	}{PublicKey: publicKey}
	out := new(SSHKey)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/users/me/ssh_keys", nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single SSH key by ID.
// Operation: getSshKey (GET /api/v1/users/me/ssh_keys/{id})
func (s *SSHKeysService) Get(ctx context.Context, id int64) (*SSHKey, error) {
	out := new(SSHKey)
	path := fmt.Sprintf("/api/v1/users/me/ssh_keys/%d", id)
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes an SSH key.
// Operation: deleteSshKey (DELETE /api/v1/users/me/ssh_keys/{id})
func (s *SSHKeysService) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/users/me/ssh_keys/%d", id)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
