package lyngbak

import (
	"io"
	"time"
)

type StaticSite struct {
	UUID             string                 `json:"uuid"`
	Name             string                 `json:"name"`
	Label            string                 `json:"label"`
	State            string                 `json:"state"`
	URL              *string                `json:"url"`
	Region           *string                `json:"region"`
	ProjectID        *string                `json:"project_id"`
	DeploymentConfig *StaticSiteConfigState `json:"deployment_config"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type StaticSiteConfigState struct {
	SourceType        string  `json:"source_type"`
	Generator         *string `json:"generator"`
	BuildCommand      *string `json:"build_command"`
	OutputDir         *string `json:"output_dir"`
	ProjectPath       *string `json:"project_path"`
	SPAMode           bool    `json:"spa_mode"`
	Repository        *string `json:"repository"`
	Branch            *string `json:"branch"`
	AutoDeployEnabled bool    `json:"auto_deploy_enabled"`
	AutoDeployBranch  *string `json:"auto_deploy_branch"`
	CredentialID      *string `json:"credential_id"`
	BucketName        *string `json:"bucket_name"`
	GitSSHURL         *string `json:"git_ssh_url"`
	SFTPUsername      *string `json:"sftp_username"`
	SFTPEndpoint      *string `json:"sftp_endpoint"`
}

type StaticSiteDeploymentConfigRequest struct {
	SourceType        string `json:"source_type"`
	CredentialID      string `json:"credential_id,omitempty"`
	Repository        string `json:"repository,omitempty"`
	Branch            string `json:"branch,omitempty"`
	AutoDeployEnabled *bool  `json:"auto_deploy_enabled,omitempty"`
	AutoDeployBranch  string `json:"auto_deploy_branch,omitempty"`
	Generator         string `json:"generator,omitempty"`
	BuildCommand      string `json:"build_command,omitempty"`
	OutputDir         string `json:"output_dir,omitempty"`
	ProjectPath       string `json:"project_path,omitempty"`
	SPAMode           *bool  `json:"spa_mode,omitempty"`
}

type StaticSiteBuildVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateStaticSiteRequest struct {
	Label                 string                            `json:"label"`
	Name                  string                            `json:"name"`
	ProjectID             string                            `json:"project_id,omitempty"`
	NewProjectName        string                            `json:"new_project_name,omitempty"`
	NewProjectDescription string                            `json:"new_project_description,omitempty"`
	Region                string                            `json:"region,omitempty"`
	DeploymentConfig      StaticSiteDeploymentConfigRequest `json:"deployment_config"`
	AppVarsAttributes     []StaticSiteBuildVar              `json:"app_vars_attributes,omitempty"`
}

type UpdateStaticSiteRequest struct {
	Label     string `json:"label,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

type UpdateStaticSiteDeploymentRequest struct {
	CredentialID      string `json:"credential_id,omitempty"`
	Repository        string `json:"repository,omitempty"`
	Branch            string `json:"branch,omitempty"`
	AutoDeployEnabled *bool  `json:"auto_deploy_enabled,omitempty"`
	AutoDeployBranch  string `json:"auto_deploy_branch,omitempty"`
	Generator         string `json:"generator,omitempty"`
	BuildCommand      string `json:"build_command,omitempty"`
	OutputDir         string `json:"output_dir,omitempty"`
	ProjectPath       string `json:"project_path,omitempty"`
	SPAMode           *bool  `json:"spa_mode,omitempty"`
}

type DeployStaticSiteRequest struct {
	Archive  io.Reader
	Filename string
}

type PreviewEnvironmentConfig struct {
	Enabled               bool           `json:"enabled"`
	ResourceMode          string         `json:"resource_mode"`
	ResourceID            *string        `json:"resource_id"`
	ProjectID             *string        `json:"project_id"`
	TriggerMode           string         `json:"trigger_mode"`
	TriggerFilter         string         `json:"trigger_filter"`
	TriggerLabel          *string        `json:"trigger_label"`
	BranchPattern         *string        `json:"branch_pattern"`
	PRNamePattern         *string        `json:"pr_name_pattern"`
	AutoDeployOnPush      bool           `json:"auto_deploy_on_push"`
	CleanupOnPRMerge      bool           `json:"cleanup_on_pr_merge"`
	CleanupOnPRClose      bool           `json:"cleanup_on_pr_close"`
	CleanupOnBranchDelete bool           `json:"cleanup_on_branch_delete"`
	CleanupOnInactivity   bool           `json:"cleanup_on_inactivity"`
	InactivityDays        int64          `json:"inactivity_days"`
	RetentionDays         int64          `json:"retention_days"`
	RetentionMode         string         `json:"retention_mode"`
	CommentOnGithub       bool           `json:"comment_on_github"`
	CloneSettings         map[string]any `json:"clone_settings"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

type UpdatePreviewEnvironmentConfigRequest struct {
	Enabled               *bool          `json:"enabled,omitempty"`
	ResourceMode          string         `json:"resource_mode,omitempty"`
	ResourceID            *string        `json:"resource_id,omitempty"`
	ProjectID             *string        `json:"project_id,omitempty"`
	TriggerMode           string         `json:"trigger_mode,omitempty"`
	TriggerFilter         string         `json:"trigger_filter,omitempty"`
	TriggerLabel          string         `json:"trigger_label,omitempty"`
	BranchPattern         string         `json:"branch_pattern,omitempty"`
	PRNamePattern         string         `json:"pr_name_pattern,omitempty"`
	AutoDeployOnPush      *bool          `json:"auto_deploy_on_push,omitempty"`
	CleanupOnPRMerge      *bool          `json:"cleanup_on_pr_merge,omitempty"`
	CleanupOnPRClose      *bool          `json:"cleanup_on_pr_close,omitempty"`
	CleanupOnBranchDelete *bool          `json:"cleanup_on_branch_delete,omitempty"`
	CleanupOnInactivity   *bool          `json:"cleanup_on_inactivity,omitempty"`
	InactivityDays        *int64         `json:"inactivity_days,omitempty"`
	RetentionDays         *int64         `json:"retention_days,omitempty"`
	RetentionMode         string         `json:"retention_mode,omitempty"`
	CommentOnGithub       *bool          `json:"comment_on_github,omitempty"`
	CloneSettings         map[string]any `json:"clone_settings,omitempty"`
}

type PreviewEnvironment struct {
	ID                int64      `json:"id"`
	Status            string     `json:"status"`
	AppUUID           *string    `json:"app_uuid"`
	AppName           *string    `json:"app_name"`
	URL               *string    `json:"url"`
	PullRequestNumber *int64     `json:"pull_request_number"`
	PullRequestTitle  *string    `json:"pull_request_title"`
	PullRequestAuthor *string    `json:"pull_request_author"`
	PullRequestState  *string    `json:"pull_request_state"`
	PullRequestURL    *string    `json:"pull_request_url"`
	PullRequestMerged *bool      `json:"pull_request_merged"`
	Branch            *string    `json:"branch"`
	BaseBranch        *string    `json:"base_branch"`
	CommitSHA         *string    `json:"commit_sha"`
	DeploymentError   *string    `json:"deployment_error"`
	LastCommitAt      *time.Time `json:"last_commit_at"`
	LastDeployedAt    *time.Time `json:"last_deployed_at"`
	ClosedAt          *time.Time `json:"closed_at"`
	ExpiresAt         *time.Time `json:"expires_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
