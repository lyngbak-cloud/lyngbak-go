package lyngbak

import "time"

type SecretFile struct {
	Name      string    `json:"name"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SecretFileRequest struct {
	Filename string `json:"filename,omitempty"`
	Text     string `json:"text,omitempty"`
}

type GitCredential struct {
	UUID           string    `json:"uuid"`
	Name           string    `json:"name"`
	Provider       string    `json:"provider"`
	InstallationID string    `json:"installation_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
