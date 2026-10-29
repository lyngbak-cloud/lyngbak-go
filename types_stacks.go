package lyngbak

import (
	"encoding/json"
	"time"
)

type Stack struct {
	UUID             string            `json:"uuid"`
	Name             string            `json:"name"`
	Label            string            `json:"label"`
	State            string            `json:"state"`
	ProjectID        *string           `json:"project_id"`
	ComposePath      *string           `json:"compose_path"`
	DeploymentConfig *StackConfigState `json:"deployment_config"`
	TotalCPUSize     float64           `json:"total_cpu_size"`
	TotalRAMSize     int64             `json:"total_ram_size"`
	TotalDiskSize    int64             `json:"total_disk_size"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Miget            *Resource         `json:"miget"`
	Services         []StackItem       `json:"services"`
	LatestDeployment *StackDeployment  `json:"latest_deployment"`
}

type StackConfigState struct {
	Type              string  `json:"type"`
	Repository        *string `json:"repository"`
	RepositoryURL     string  `json:"repository_url"`
	Branch            *string `json:"branch"`
	AutoDeployEnabled *bool   `json:"auto_deploy_enabled"`
	CredentialID      *string `json:"credential_id"`
}

type StackItem struct {
	ComposeServiceName string  `json:"compose_service_name"`
	ComposeServiceType string  `json:"compose_service_type"`
	Status             string  `json:"status"`
	ItemType           string  `json:"item_type"`
	ItemUUID           *string `json:"item_uuid"`
	ItemName           *string `json:"item_name"`
}

type StackDeployment struct {
	UUID             string          `json:"uuid"`
	State            string          `json:"state"`
	Subevent         string          `json:"subevent"`
	CommitSHA        *string         `json:"commit_sha"`
	ValidationErrors json.RawMessage `json:"validation_errors"`
	PublishedAt      *time.Time      `json:"published_at"`
	CompletedAt      *time.Time      `json:"completed_at"`
	CreatedAt        time.Time       `json:"created_at"`
}

type StackAnalysis struct {
	Manifest map[string]any `json:"manifest"`
	Warnings []any          `json:"warnings"`
}

type CreateStackRequest struct {
	RepositoryURL            string                       `json:"repository_url"`
	Branch                   string                       `json:"branch"`
	ResourceID               string                       `json:"resource_id"`
	ComposePath              string                       `json:"compose_path,omitempty"`
	CredentialID             string                       `json:"credential_id,omitempty"`
	Label                    string                       `json:"label,omitempty"`
	Name                     string                       `json:"name,omitempty"`
	ProjectID                string                       `json:"project_id,omitempty"`
	NewProjectName           string                       `json:"new_project_name,omitempty"`
	NewProjectDescription    string                       `json:"new_project_description,omitempty"`
	EnvVarOverrides          map[string]map[string]string `json:"env_var_overrides,omitempty"`
	AutoPopulateRequiredVars bool                         `json:"auto_populate_required_vars,omitempty"`
}

type AnalyzeStackRequest struct {
	RepositoryURL string `json:"repository_url"`
	Branch        string `json:"branch"`
	ComposePath   string `json:"compose_path,omitempty"`
	ResourceID    string `json:"resource_id,omitempty"`
	CredentialID  string `json:"credential_id,omitempty"`
}

type UpdateStackRequest struct {
	Label       *string `json:"label,omitempty"`
	ComposePath *string `json:"compose_path,omitempty"`
	ProjectID   *string `json:"project_id,omitempty"`
}

type DeployStackRequest struct {
	CommitSHA string `json:"commit_sha,omitempty"`
}

type UpdateStackDeploymentRequest struct {
	Branch            *string `json:"branch,omitempty"`
	AutoDeployEnabled *bool   `json:"auto_deploy_enabled,omitempty"`
	Repository        *string `json:"repository,omitempty"`
}
