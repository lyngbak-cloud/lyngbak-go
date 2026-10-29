package lyngbak

import "time"

// This file contains all request and response types, derived from the
// components.schemas section of the OpenAPI specification
// (openapi/openapi.json). Entity (response) types use value fields because
// the spec marks all entity properties as required; fields the spec leaves
// loosely typed are noted in comments.

// ----------------------------------------------------------------------------
// Entities (response types)
// ----------------------------------------------------------------------------

// Message is a generic API response message.
// Schema: Api_V1_Entities_Message.
type Message struct {
	Message string `json:"message"`
}

// App is a Miget application.
// Schema: Api_V1_Entities_App.
type App struct {
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	State            string `json:"state"`
	ProjectID        string `json:"project_id"`
	DeploymentMethod string `json:"deployment_method"`
	// DeploymentConfig fields vary by deployment method; the spec declares
	// a free-form object.
	DeploymentConfig *AppDeploymentConfigState `json:"deployment_config"`
	Builder          string                    `json:"builder"`
	ManageDNSRecords bool                      `json:"manage_dns_records"`
	Maintenance      bool                      `json:"maintenance"`
	ServiceName      string                    `json:"service_name"`
	AllowConnections bool                      `json:"allow_connections"`
	// Quota is the resource quota configuration; the spec declares a
	// free-form object.
	Quota     map[string]any `json:"quota"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Miget     Resource       `json:"miget"`
	Region    Region         `json:"region"`
	Addons    []Addon        `json:"addons"`
	// Cronjobs holds the app's scheduled cron jobs, rendered by the server
	// as Addon entities; the schedule is in Settings.
	Cronjobs []Addon `json:"cronjobs"`

	HealthChecks              *AppHealthChecks   `json:"health_checks"`
	ScalingProfile            *AppScalingProfile `json:"scaling_profile"`
	BasicAuthEnabled          bool               `json:"basic_auth_enabled"`
	BasicAuthUsername         *string            `json:"basic_auth_username"`
	ProjectVariablesEnabled   bool               `json:"project_variables_enabled"`
	ProjectSecretFilesEnabled bool               `json:"project_secret_files_enabled"`
	RAMSize                   float64            `json:"ram_size"`
	CPUSize                   float64            `json:"cpu_size"`
	ResourceID                string             `json:"resource_id"`
	VPCUUID                   *string            `json:"vpc_uuid"`
	VPCSubnetUUID             *string            `json:"vpc_subnet_uuid"`
	PrivateAccess             bool               `json:"private_access"`
	InternalURL               *string            `json:"internal_url"`
	PublicURL                 *string            `json:"public_url"`
}

type AppHealthChecks struct {
	LivenessProbeEnabled                     *bool   `json:"liveness_probe_enabled"`
	ReadinessProbeEnabled                    *bool   `json:"readiness_probe_enabled"`
	StartupProbeEnabled                      *bool   `json:"startup_probe_enabled"`
	LivenessInAppFailureNotificationEnabled  *bool   `json:"liveness_in_app_failure_notification_enabled"`
	LivenessEmailFailureNotificationEnabled  *bool   `json:"liveness_email_failure_notification_enabled"`
	ReadinessInAppFailureNotificationEnabled *bool   `json:"readiness_in_app_failure_notification_enabled"`
	ReadinessEmailFailureNotificationEnabled *bool   `json:"readiness_email_failure_notification_enabled"`
	StartupInAppFailureNotificationEnabled   *bool   `json:"startup_in_app_failure_notification_enabled"`
	StartupEmailFailureNotificationEnabled   *bool   `json:"startup_email_failure_notification_enabled"`
	LivenessProbePath                        *string `json:"liveness_probe_path"`
	ReadinessProbePath                       *string `json:"readiness_probe_path"`
	StartupProbePath                         *string `json:"startup_probe_path"`
	LivenessProbeInitialDelaySeconds         *int    `json:"liveness_probe_initial_delay_seconds"`
	LivenessProbeTimeoutSeconds              *int    `json:"liveness_probe_timeout_seconds"`
	LivenessProbePeriodSeconds               *int    `json:"liveness_probe_period_seconds"`
	LivenessProbeFailureThreshold            *int    `json:"liveness_probe_failure_threshold"`
	ReadinessProbeInitialDelaySeconds        *int    `json:"readiness_probe_initial_delay_seconds"`
	ReadinessProbeTimeoutSeconds             *int    `json:"readiness_probe_timeout_seconds"`
	ReadinessProbePeriodSeconds              *int    `json:"readiness_probe_period_seconds"`
	ReadinessProbeFailureThreshold           *int    `json:"readiness_probe_failure_threshold"`
	StartupProbeInitialDelaySeconds          *int    `json:"startup_probe_initial_delay_seconds"`
	StartupProbeTimeoutSeconds               *int    `json:"startup_probe_timeout_seconds"`
	StartupProbePeriodSeconds                *int    `json:"startup_probe_period_seconds"`
	StartupProbeFailureThreshold             *int    `json:"startup_probe_failure_threshold"`
}

type AppScalingProfile struct {
	Replicas           int     `json:"replicas"`
	AutoScalingEnabled bool    `json:"auto_scaling_enabled"`
	AutoMinReplicas    *int    `json:"auto_min_replicas"`
	AutoMaxReplicas    *int    `json:"auto_max_replicas"`
	ScalingStartTime   *string `json:"scaling_start_time"`
	ScalingEndTime     *string `json:"scaling_end_time"`
	CPUThreshold       *int    `json:"cpu_threshold"`
	MemoryThreshold    *int    `json:"memory_threshold"`
}

type AppDeploymentConfigState struct {
	Type                string   `json:"type"`
	Repository          *string  `json:"repository"`
	RepositoryURL       *string  `json:"repository_url"`
	Branch              *string  `json:"branch"`
	AutoDeployEnabled   *bool    `json:"auto_deploy_enabled"`
	AutoDeployBranch    *string  `json:"auto_deploy_branch"`
	CredentialID        *string  `json:"credential_id"`
	DockerfilePath      *string  `json:"dockerfile_path"`
	BuildContext        *string  `json:"build_context"`
	ProjectPath         *string  `json:"project_path"`
	RunCommand          *string  `json:"run_command"`
	Language            *string  `json:"language"`
	BuildCommand        *string  `json:"build_command"`
	PreDeployCommand    *string  `json:"pre_deploy_command"`
	PostDeployCommand   *string  `json:"post_deploy_command"`
	UseDHI              *bool    `json:"use_dhi"`
	ImageURL            *string  `json:"image_url"`
	Tag                 *string  `json:"tag"`
	Command             []string `json:"command"`
	Args                []string `json:"args"`
	ParentAppID         *string  `json:"parent_app_id"`
	ParentImageAutoSync *bool    `json:"parent_image_auto_sync"`
	SSHKeys             []string `json:"ssh_keys"`
	RegistryImage       *string  `json:"registry_image"`
	RegistryHostname    *string  `json:"registry_hostname"`
	RegistryUsername    *string  `json:"registry_username"`
	RegistryPassword    *string  `json:"registry_password"`
}

// Resource is a Miget compute resource.
// Schema: Api_V1_Entities_Resource.
type Resource struct {
	UUID              string      `json:"uuid"`
	Name              string      `json:"name"`
	State             string      `json:"state"`
	Labels            []string    `json:"labels"`
	TotalDiskSize     float64     `json:"total_disk_size"`      // GiB
	TotalRAMSize      float64     `json:"total_ram_size"`       // MiB
	TotalCPUSize      float64     `json:"total_cpu_size"`       // cores
	TotalUsedDiskSize float64     `json:"total_used_disk_size"` // GiB
	TotalUsedRAMSize  float64     `json:"total_used_ram_size"`  // MiB
	TotalUsedCPUSize  float64     `json:"total_used_cpu_size"`  // cores
	AvailableCPUSize  float64     `json:"available_cpu_size"`   // cores
	AvailableDiskSize float64     `json:"available_disk_size"`  // GiB
	AvailableRAMSize  float64     `json:"available_ram_size"`   // MiB
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Plan              Plan        `json:"plan"`
	Components        []Component `json:"components"`

	RegionCode          string              `json:"region_code"`
	ComponentQuantities []ComponentQuantity `json:"component_quantities"`
	Assigned            bool                `json:"assigned"`
	ProjectIDs          []string            `json:"project_ids"`
}

type ComponentQuantity struct {
	CodeName string `json:"code_name"`
	Quantity int64  `json:"quantity"`
}

// Plan is a billing plan.
// Schema: Api_V1_Entities_Plan.
type Plan struct {
	Name        string  `json:"name"`
	CodeName    string  `json:"code_name"`
	Description string  `json:"description"`
	PlanType    string  `json:"plan_type"`
	RAMSize     float64 `json:"ram_size"`   // MiB
	DiskSize    float64 `json:"disk_size"`  // GiB
	CPUSize     float64 `json:"cpu_size"`   // cores
	UnitPrice   int64   `json:"unit_price"` // cents per billing period
}

// Component is an additional resource component (extra RAM, CPU, disk).
// Schema: Api_V1_Entities_Component.
type Component struct {
	Name      string `json:"name"`
	CodeName  string `json:"code_name"`
	UnitPrice int64  `json:"unit_price"` // cents per unit
}

// Region is a deployment region.
// Schema: Api_V1_Entities_Region.
type Region struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// Addon is an application add-on (postgres, mysql, valkey, storage, cronjob).
// Schema: Api_V1_Entities_Addon.
type Addon struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
	State string `json:"state"`
	// Quota, Settings, Var and ConnectionDetails are declared as free-form
	// objects in the spec.
	Quota             map[string]any `json:"quota"`
	Settings          map[string]any `json:"settings"`
	Var               map[string]any `json:"var"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	ConnectionDetails map[string]any `json:"connection_details"`

	RAMSize          float64        `json:"ram_size"`
	DiskSize         float64        `json:"disk_size"`
	CPUSize          float64        `json:"cpu_size"`
	ServiceID        *int64         `json:"service_id"`
	BackupEnabled    *bool          `json:"backup_enabled"`
	ScheduledBackup  *string        `json:"scheduled_backup"`
	Role             *string        `json:"role"`
	PrimaryAddonUUID *string        `json:"primary_addon_uuid"`
	Replicas         []AddonReplica `json:"replicas"`
}

type AddonReplica struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Label string `json:"label"`
	State string `json:"state"`
}

// Activity is an audit log entry for an app.
// Schema: Api_V1_Entities_Activity.
type Activity struct {
	Action      string `json:"action"`
	Description string `json:"description"`
	// Resource describes the resource the activity is about; the spec
	// documents fields type, name and uuid but declares a free-form object.
	Resource map[string]any `json:"resource"`
	// Actor describes the user who performed the action (fields uuid, name,
	// email); null for system actions.
	Actor     map[string]any `json:"actor"`
	Timestamp time.Time      `json:"timestamp"`
	Source    string         `json:"source"` // api, ui
}

// Domain is a custom domain attached to an app.
// Schema: Api_V1_Entities_Domain.
type Domain struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	// DNSTarget is populated only after verification succeeds.
	DNSTarget string `json:"dns_target"`
	// VerificationStatus is one of pending, in_progress, verified, failed.
	VerificationStatus string `json:"verification_status"`
	// VerificationToken is the token to publish as a TXT record at
	// _migetapp-verify.<domain>. Populated only while verification is pending.
	VerificationToken string    `json:"verification_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Backup describes the backup configuration and history of a database
// add-on or service.
// Schema: Api_V1_Entities_Backup.
type Backup struct {
	BackupEnabled bool `json:"backup_enabled"`
	BackupActive  bool `json:"backup_active"`
	// Backups, ScheduledBackup and ClusterStatus are declared as free-form
	// objects in the spec.
	Backups         []map[string]any `json:"backups"`
	ScheduledBackup map[string]any   `json:"scheduled_backup"`
	ClusterStatus   map[string]any   `json:"cluster_status"`
}

// Variable is an environment variable (app- or project-level).
// Schema: Api_V1_Entities_Variable.
type Variable struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Cronjob is a scheduled cron job attached to an app.
// Schema: Api_V1_Entities_Cronjob.
type Cronjob struct {
	UUID           string    `json:"uuid"`
	Name           string    `json:"name"`
	ScheduleType   string    `json:"schedule_type"` // cron, interval
	IntervalType   string    `json:"interval_type"` // every_10_minutes, hourly, daily
	Hour           string    `json:"hour"`
	Minute         string    `json:"minute"`
	Command        string    `json:"command"`
	Label          string    `json:"label"`
	CronExpression string    `json:"cron_expression"`
	Status         string    `json:"status"` // running, stopped
	LastJobName    string    `json:"last_job_name"`
	LastJobStatus  string    `json:"last_job_status"`
	LastJobRunAt   time.Time `json:"last_job_run_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Port is an exposed application port.
// Schema: Api_V1_Entities_Port.
type Port struct {
	ID                    int64  `json:"id"`
	InternalPort          int64  `json:"internal_port"`
	ExternalPort          *int64 `json:"external_port"`
	Protocol              string `json:"protocol"`
	AutoDetected          bool   `json:"auto_detected"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	Public                bool   `json:"public"`
	PublicEndpoint        string `json:"public_endpoint"`
	EndpointID            int64  `json:"endpoint_id"`
	EndpointHostname      string `json:"endpoint_hostname"`
	EndpointIPv4Primary   string `json:"endpoint_ipv4_primary"`
	EndpointIPv4Secondary string `json:"endpoint_ipv4_secondary"`
	EndpointIPv6Primary   string `json:"endpoint_ipv6_primary"`
	EndpointIPv6Secondary string `json:"endpoint_ipv6_secondary"`
}

// Deployment is an application deployment.
// Schema: Api_V1_Entities_Deployment.
type Deployment struct {
	UUID string `json:"uuid"`
	// State is one of pending, running, completed, failed, cancelling,
	// cancelled.
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// LogsStoredAt is the zero time if logs have not been stored yet
	// (the API returns null).
	LogsStoredAt  time.Time `json:"logs_stored_at"`
	Branch        *string   `json:"branch"`
	CommitSHA     *string   `json:"commit_sha"`
	CommitMessage *string   `json:"commit_message"`
}

// AuthTokens is the result of signing in.
// Schema: Api_V1_Entities_AuthTokens.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`  // JWT access token (30 min expiry)
	RefreshToken string `json:"refresh_token"` // JWT refresh token (7 days expiry)
}

// AccessToken is the result of refreshing an access token.
// Schema: Api_V1_Entities_AccessToken.
type AccessToken struct {
	AccessToken string `json:"access_token"`
}

// Bucket is an S3-compatible object storage bucket.
// Schema: Api_V1_Entities_Bucket.
type Bucket struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	State       string `json:"state"`      // pending, active, failed, blocked
	Visibility  string `json:"visibility"` // public_access, private_access
	S3Endpoint  string `json:"s3_endpoint"`
	S3AccessKey string `json:"s3_access_key"`
	S3SecretKey string `json:"s3_secret_key"`
	// PublicEndpoint is set only for public buckets.
	PublicEndpoint string `json:"public_endpoint"`
	RegionCode     string `json:"region_code"`
	// UsageBytes and QuotaBytes are declared int32 in the spec but are byte
	// counts; int64 is used to avoid overflow.
	UsageBytes      int64     `json:"usage_bytes"`
	QuotaBytes      int64     `json:"quota_bytes"`
	UsagePercentage float64   `json:"usage_percentage"`
	ObjectCount     int64     `json:"object_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Policy is the S3-compatible bucket policy as a JSON string.
	Policy string `json:"policy"`
	// ACL is the S3-compatible access control list.
	ACL        string   `json:"acl"`
	Miget      Resource `json:"miget"`
	ResourceID string   `json:"resource_id"`
	ProjectID  *string  `json:"project_id"`
}

// BucketObjectList is a page of bucket objects.
// Schema: Api_V1_Entities_BucketObjectList.
type BucketObjectList struct {
	Objects     []BucketObject `json:"objects"`
	CurrentPath string         `json:"current_path"`
	HasMore     bool           `json:"has_more"`
	// NextCursor is the pagination cursor for the next page; empty on the
	// last page.
	NextCursor string `json:"next_cursor"`
}

// BucketObject is a file or folder inside a bucket.
type BucketObject struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	IsFolder bool   `json:"is_folder"`
	// Size is declared int32 in the spec but is a byte count; int64 is used
	// to avoid overflow.
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	ContentType  string    `json:"content_type"`
}

// ContainerRegistryCredential is a stored container registry credential.
// Schema: Api_V1_Entities_ContainerRegistryCredential.
type ContainerRegistryCredential struct {
	// UUID is passed as deployment_config.credential_id when creating apps.
	UUID string `json:"uuid"`
	Name string `json:"name"`
	// Registry is one of docker_hub, github, gitlab, aws_ecr, azure,
	// digitalocean, quay, generic.
	Registry         string    `json:"registry"`
	Username         string    `json:"username"`
	RegistryHostname string    `json:"registry_hostname"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Project is a project grouping apps.
// Schema: Api_V1_Entities_Project.
type Project struct {
	UUID         string               `json:"uuid"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Resources    []Resource           `json:"resources"`
	Restrictions []ProjectRestriction `json:"restrictions"`
}

type ProjectRestriction struct {
	ID          int64   `json:"id"`
	SubjectType string  `json:"subject_type"`
	UserUUID    *string `json:"user_uuid"`
	RoleName    *string `json:"role_name"`
}

// Service is a standalone service (postgres, shared_storage).
// Schema: Api_V1_Entities_Service.
type Service struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	ServiceType string `json:"service_type"` // postgres, shared_storage
	ProjectID   string `json:"project_id"`
	ResourceID  string `json:"resource_id"`
	// MigetID is deprecated; use ResourceID instead.
	MigetID   *string   `json:"miget_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// ConnectionDetails is declared as a free-form object in the spec.
	ConnectionDetails map[string]any `json:"connection_details"`

	RAMSize         float64        `json:"ram_size"`
	DiskSize        float64        `json:"disk_size"`
	CPUSize         float64        `json:"cpu_size"`
	PublicAccess    *bool          `json:"public_access"`
	Instances       *int           `json:"instances"`
	PostgresVersion *string        `json:"postgres_version"`
	MountPoint      *string        `json:"mount_point"`
	Mounts          []ServiceMount `json:"mounts"`
	BackupEnabled   *bool          `json:"backup_enabled"`
	ScheduledBackup *string        `json:"scheduled_backup"`
	State           string         `json:"state"`

	Role             string         `json:"role"`
	PrimaryAddonUUID *string        `json:"primary_addon_uuid"`
	Replicas         []Service      `json:"replicas"`
	Addons           []ServiceAddon `json:"addons"`
}

type ServiceAddon struct {
	UUID      string         `json:"uuid"`
	Name      string         `json:"name"`
	Label     string         `json:"label"`
	State     string         `json:"state"`
	Settings  map[string]any `json:"settings"`
	IsReplica bool           `json:"is_replica"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ServiceMount struct {
	AppID      string  `json:"app_id"`
	AddonUUID  string  `json:"addon_uuid"`
	MountPoint string  `json:"mount_point"`
	SubPath    *string `json:"sub_path"`
	Label      *string `json:"label"`
}

// User is the authenticated user.
// Schema: Api_V1_Entities_User.
type User struct {
	UUID            string    `json:"uuid"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`    // first name
	Surname         string    `json:"surname"` // last name
	DiscordUsername string    `json:"discord_username"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SSHKey is a user SSH public key.
// Schema: Api_V1_Entities_SshKey.
type SSHKey struct {
	ID        int64     `json:"id"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Credit is the user's credit balance and history.
// Schema: Api_V1_Entities_Credit.
type Credit struct {
	CreditBalanceInCents   int64  `json:"credit_balance_in_cents"`
	CreditBalanceFormatted string `json:"credit_balance_formatted"`
	ReferralCode           string `json:"referral_code"`
	ReferralPercentage     int64  `json:"referral_percentage"`
	// CreditOperations items are not described by the spec.
	CreditOperations []any `json:"credit_operations"`
}

// ----------------------------------------------------------------------------
// Enum-like constants
// ----------------------------------------------------------------------------

// App state actions accepted by AppsService.UpdateState.
const (
	AppStateScheduleStart   = "schedule_start"
	AppStateScheduleStop    = "schedule_stop"
	AppStateScheduleRestart = "schedule_restart"
)

// Process state actions accepted by AddonsService.UpdateState and
// ServicesService.UpdateState.
const (
	ProcessStateStart   = "process_start"
	ProcessStateStop    = "process_stop"
	ProcessStateRestart = "process_restart"
)

// ----------------------------------------------------------------------------
// Request types
// ----------------------------------------------------------------------------

// AppVarAttribute is an environment variable set at app creation time.
type AppVarAttribute struct {
	Key   string `json:"key"` // SCREAMING_SNAKE_CASE
	Value string `json:"value"`
}

// AppDeploymentConfig is the deployment configuration sent when creating an
// app. The relevant fields depend on the deployment method.
type AppDeploymentConfig struct {
	// CredentialID is the UUID of stored credentials (github, public_git,
	// container_registry).
	CredentialID string `json:"credential_id,omitempty"`
	// ImageURL is the container image URL (container_registry).
	ImageURL string `json:"image_url,omitempty"`
	// Tag is the container image tag (container_registry).
	Tag string `json:"tag,omitempty"`
	// Repository is the repository name (github).
	Repository string `json:"repository,omitempty"`
	// Branch is the branch name (github, public_git).
	Branch string `json:"branch,omitempty"`
	// DockerfilePath is the path to the Dockerfile (github, public_git,
	// git_push).
	DockerfilePath string `json:"dockerfile_path,omitempty"`
	// BuildContext is the Docker build context path (github, public_git,
	// git_push).
	BuildContext string `json:"build_context,omitempty"`
	// AutoDeployEnabled enables auto-deploy on push (github).
	AutoDeployEnabled *bool `json:"auto_deploy_enabled,omitempty"`
	// AutoDeployBranch is the branch for auto-deploy (github).
	AutoDeployBranch string `json:"auto_deploy_branch,omitempty"`
	// ParentAppID is the UUID of the parent app (parent_image).
	ParentAppID string `json:"parent_app_id,omitempty"`
	// ParentImageAutoSync auto-syncs with the parent image (parent_image).
	ParentImageAutoSync *bool `json:"parent_image_auto_sync,omitempty"`
	// SSHKeys are SSH public keys for deployment access (kamal).
	SSHKeys []string `json:"ssh_keys,omitempty"`
	// Command overrides the container ENTRYPOINT (container_registry).
	Command []string `json:"command,omitempty"`
	// Args overrides the container CMD (container_registry).
	Args []string `json:"args,omitempty"`
}

// CreateAppRequest is the body for AppsService.Create.
// Schema: createApp.
type CreateAppRequest struct {
	PrivateAccess *bool  `json:"private_access,omitempty"`
	VPCUUID       string `json:"vpc_uuid,omitempty"`
	VPCSubnetUUID string `json:"vpc_subnet_uuid,omitempty"`
	Label         string `json:"label"`
	Name          string `json:"name"`
	ProjectID     string `json:"project_id"`
	ResourceID    string `json:"resource_id"`
	// Builder is one of auto, dockerfile, custom.
	Builder string   `json:"builder"`
	RAMSize *float64 `json:"ram_size,omitempty"` // MiB
	CPUSize *float64 `json:"cpu_size,omitempty"` // cores
	// DeploymentMethod is one of git_push, parent_image, container_registry,
	// github, public_git, kamal.
	DeploymentMethod string               `json:"deployment_method,omitempty"`
	DeploymentConfig *AppDeploymentConfig `json:"deployment_config,omitempty"`
	AppVars          []AppVarAttribute    `json:"app_vars_attributes,omitempty"`
}

// UpdateAppRequest is the body for AppsService.Update.
// Schema: updateApp.
type UpdateAppRequest struct {
	PrivateAccess *bool  `json:"private_access,omitempty"`
	Label         string `json:"label,omitempty"`
	ProjectID     string `json:"project_id,omitempty"`
	// Builder is one of auto, dockerfile, custom.
	Builder          string   `json:"builder,omitempty"`
	ManageDNSRecords *bool    `json:"manage_dns_records,omitempty"`
	Maintenance      *bool    `json:"maintenance,omitempty"`
	RAMSize          *float64 `json:"ram_size,omitempty"` // MiB
	CPUSize          *float64 `json:"cpu_size,omitempty"` // cores
}

// UpdateAppSecurityRequest is the body for AppsService.UpdateSecurity.
// Schema: updateAppSecurity.
type UpdateAppSecurityRequest struct {
	AllowConnections *bool `json:"allow_connections,omitempty"`
	BasicAuthEnabled *bool `json:"basic_auth_enabled,omitempty"`
	// BasicAuthUsername is required when BasicAuthEnabled is true.
	BasicAuthUsername string `json:"basic_auth_username,omitempty"`
	// BasicAuthPassword is required when BasicAuthEnabled is true; leave
	// empty to keep the current password.
	BasicAuthPassword string `json:"basic_auth_password,omitempty"`
}

// CloneAppAddon selects an add-on to clone with CloneAppRequest.
type CloneAppAddon struct {
	UUID string `json:"uuid"`
	// CloneData includes the add-on data (database contents, files).
	CloneData bool `json:"clone_data"`
}

// CloneAppRequest is the body for AppsService.Clone.
// Schema: cloneApp.
type CloneAppRequest struct {
	Label                string          `json:"label"`
	Name                 string          `json:"name"`
	ProjectID            string          `json:"project_id"`
	ResourceID           string          `json:"resource_id"`
	UseParentImage       *bool           `json:"use_parent_image,omitempty"`
	ParentImageAutoSync  *bool           `json:"parent_image_auto_sync,omitempty"`
	CloneVariables       *bool           `json:"clone_variables,omitempty"`
	CloneSecretFiles     *bool           `json:"clone_secret_files,omitempty"`
	CloneScalingSettings *bool           `json:"clone_scaling_settings,omitempty"`
	CloneHealthChecks    *bool           `json:"clone_health_checks,omitempty"`
	CloneSecurity        *bool           `json:"clone_security,omitempty"`
	Addons               []CloneAppAddon `json:"addons,omitempty"`
	// Cronjobs is a list of cron job UUIDs to clone.
	Cronjobs []string `json:"cronjobs,omitempty"`
}

// DeployAppRequest is the body for AppsService.Deploy. All fields are
// optional; a nil request triggers a default deployment.
// Schema: deployApp.
type DeployAppRequest struct {
	// CustomTag is a custom image tag to deploy.
	CustomTag string `json:"custom_tag,omitempty"`
	// CommitSHA pins a specific commit (Git/GitHub deployments).
	CommitSHA string `json:"commit_sha,omitempty"`
	// Branch selects the branch to deploy (GitHub deployments).
	Branch string `json:"branch,omitempty"`
}

// UpdateAppDeploymentConfig is the deployment configuration sent when
// switching an app's deployment method. The relevant fields depend on the
// deployment method.
type UpdateAppDeploymentConfig struct {
	// CredentialID is the UUID of stored container registry credentials.
	CredentialID string `json:"credential_id,omitempty"`
	// ImageURL is the full container image URL (container_registry).
	ImageURL string `json:"image_url,omitempty"`
	// Tag is the container image tag (container_registry).
	Tag string `json:"tag,omitempty"`
	// Command overrides the container ENTRYPOINT; leave unset to use the
	// image default.
	Command []string `json:"command,omitempty"`
	// Args overrides the container CMD; leave unset to use the image default.
	Args []string `json:"args,omitempty"`
	// ParentAppID is the UUID of the parent application (parent_image).
	ParentAppID string `json:"parent_app_id,omitempty"`
	// ParentImageAutoSync auto-deploys when the parent image is updated
	// (parent_image).
	ParentImageAutoSync *bool `json:"parent_image_auto_sync,omitempty"`
	// SSHKeys are SSH public keys for Kamal deployment access (kamal).
	SSHKeys []string `json:"ssh_keys,omitempty"`
}

// UpdateAppDeploymentRequest is the body for AppsService.UpdateDeployment.
// Schema: updateAppDeployment.
type UpdateAppDeploymentRequest struct {
	// DeploymentMethod is one of git_push, parent_image, container_registry,
	// github, public_git, kamal.
	DeploymentMethod string                     `json:"deployment_method"`
	DeploymentConfig *UpdateAppDeploymentConfig `json:"deployment_config_attributes,omitempty"`
}

// UpdateAppScalingProfileRequest is the body for
// AppsService.UpdateScalingProfile.
// Schema: updateAppScalingProfile.
type UpdateAppScalingProfileRequest struct {
	// Replicas is the fixed number of instances (ignored when auto-scaling
	// is enabled).
	Replicas           *int64 `json:"replicas,omitempty"`
	AutoScalingEnabled *bool  `json:"auto_scaling_enabled,omitempty"`
	AutoMinReplicas    *int64 `json:"auto_min_replicas,omitempty"`
	AutoMaxReplicas    *int64 `json:"auto_max_replicas,omitempty"`
	// ScalingStartTime and ScalingEndTime use HH:MM format (24-hour).
	ScalingStartTime string `json:"scaling_start_time,omitempty"`
	ScalingEndTime   string `json:"scaling_end_time,omitempty"`
	WithinResources  *bool  `json:"within_resources,omitempty"`
	// CPUThreshold and MemoryThreshold are usage percentages (1-100).
	CPUThreshold    *int64 `json:"cpu_threshold,omitempty"`
	MemoryThreshold *int64 `json:"memory_threshold,omitempty"`
	PeriodEnabled   *bool  `json:"period_enabled,omitempty"`
}

// UpdateAppHealthChecksRequest is the body for
// AppsService.UpdateHealthChecks.
// Schema: updateAppHealthChecks.
type UpdateAppHealthChecksRequest struct {
	LivenessProbeEnabled                     *bool  `json:"liveness_probe_enabled,omitempty"`
	ReadinessProbeEnabled                    *bool  `json:"readiness_probe_enabled,omitempty"`
	StartupProbeEnabled                      *bool  `json:"startup_probe_enabled,omitempty"`
	LivenessInAppFailureNotificationEnabled  *bool  `json:"liveness_in_app_failure_notification_enabled,omitempty"`
	LivenessEmailFailureNotificationEnabled  *bool  `json:"liveness_email_failure_notification_enabled,omitempty"`
	ReadinessInAppFailureNotificationEnabled *bool  `json:"readiness_in_app_failure_notification_enabled,omitempty"`
	ReadinessEmailFailureNotificationEnabled *bool  `json:"readiness_email_failure_notification_enabled,omitempty"`
	StartupInAppFailureNotificationEnabled   *bool  `json:"startup_in_app_failure_notification_enabled,omitempty"`
	StartupEmailFailureNotificationEnabled   *bool  `json:"startup_email_failure_notification_enabled,omitempty"`
	ProjectVariablesEnabled                  *bool  `json:"project_variables_enabled,omitempty"`
	LivenessProbePath                        string `json:"liveness_probe_path,omitempty"`
	ReadinessProbePath                       string `json:"readiness_probe_path,omitempty"`
	StartupProbePath                         string `json:"startup_probe_path,omitempty"`
	LivenessProbeInitialDelaySeconds         *int64 `json:"liveness_probe_initial_delay_seconds,omitempty"`
	LivenessProbeTimeoutSeconds              *int64 `json:"liveness_probe_timeout_seconds,omitempty"`
	LivenessProbePeriodSeconds               *int64 `json:"liveness_probe_period_seconds,omitempty"`
	LivenessProbeFailureThreshold            *int64 `json:"liveness_probe_failure_threshold,omitempty"`
	ReadinessProbeInitialDelaySeconds        *int64 `json:"readiness_probe_initial_delay_seconds,omitempty"`
	ReadinessProbeTimeoutSeconds             *int64 `json:"readiness_probe_timeout_seconds,omitempty"`
	ReadinessProbePeriodSeconds              *int64 `json:"readiness_probe_period_seconds,omitempty"`
	ReadinessProbeFailureThreshold           *int64 `json:"readiness_probe_failure_threshold,omitempty"`
	StartupProbeInitialDelaySeconds          *int64 `json:"startup_probe_initial_delay_seconds,omitempty"`
	StartupProbeTimeoutSeconds               *int64 `json:"startup_probe_timeout_seconds,omitempty"`
	StartupProbePeriodSeconds                *int64 `json:"startup_probe_period_seconds,omitempty"`
	StartupProbeFailureThreshold             *int64 `json:"startup_probe_failure_threshold,omitempty"`
}

// ActivityListOptions are the query parameters for AppsService.Activity.
type ActivityListOptions struct {
	Page  int // page number (zero value omits the parameter)
	Limit int // page size (zero value omits the parameter)
}

// CreateAppAddonRequest is the body for AddonsService.Create.
// Schema: createAppAddon.
type CreateAppAddonRequest struct {
	// Type is one of postgres, mysql, valkey, storage.
	Type     string   `json:"type"`
	Label    string   `json:"label,omitempty"`
	RAMSize  *float64 `json:"ram_size,omitempty"`  // MiB
	DiskSize *float64 `json:"disk_size,omitempty"` // GiB
	CPUSize  *float64 `json:"cpu_size,omitempty"`  // cores
	// PostgresVersion, MySQLVersion and ValkeyVersion select the engine
	// version for the respective add-on type.
	PostgresVersion string `json:"postgres_version,omitempty"`
	MySQLVersion    string `json:"mysql_version,omitempty"`
	ValkeyVersion   string `json:"valkey_version,omitempty"`
	// PublicAccess enables public internet access: "0" or "1" (postgres).
	PublicAccess string `json:"public_access,omitempty"`
	// Instances is 1 for standalone, or 3/5/7 for an HA cluster (postgres).
	Instances *int64 `json:"instances,omitempty"`
	// CreationMode is fresh (default) or external_replica.
	CreationMode string `json:"creation_mode,omitempty"`
	// External replica settings (creation_mode external_replica).
	ExternalHost string `json:"external_host,omitempty"`
	ExternalPort *int64 `json:"external_port,omitempty"`
	// AuthType is password or tls.
	AuthType            string `json:"auth_type,omitempty"`
	ReplicationUsername string `json:"replication_username,omitempty"`
	ReplicationPassword string `json:"replication_password,omitempty"`
	CACrt               string `json:"ca_crt,omitempty"`
	TLSCrt              string `json:"tls_crt,omitempty"`
	TLSKey              string `json:"tls_key,omitempty"`
	// S3Enabled enables S3 WAL archive fallback: "0" or "1".
	S3Enabled    string `json:"s3_enabled,omitempty"`
	S3Endpoint   string `json:"s3_endpoint,omitempty"`
	S3Bucket     string `json:"s3_bucket,omitempty"`
	S3Path       string `json:"s3_path,omitempty"`
	S3AccessKey  string `json:"s3_access_key,omitempty"`
	S3SecretKey  string `json:"s3_secret_key,omitempty"`
	S3ServerName string `json:"s3_server_name,omitempty"`
	// ServiceID associates a shared storage add-on with a service.
	ServiceID *int64 `json:"service_id,omitempty"`
	// MountPoint is the container mount path, e.g. /data (storage).
	MountPoint string `json:"mount_point,omitempty"`
	// StorageAccess is RWO or RWX (storage).
	StorageAccess string `json:"storage_access,omitempty"`
	SubPath       string `json:"sub_path,omitempty"`
}

// UpdateAppAddonRequest is the body for AddonsService.Update.
// Schema: updateAppAddon.
type UpdateAppAddonRequest struct {
	Label    string   `json:"label,omitempty"`
	RAMSize  *float64 `json:"ram_size,omitempty"`  // MiB
	DiskSize *float64 `json:"disk_size,omitempty"` // GiB
	CPUSize  *float64 `json:"cpu_size,omitempty"`  // cores
	// PublicAccess enables public internet access: "0" or "1" (postgres).
	PublicAccess string `json:"public_access,omitempty"`
	// BackupEnabled enables automatic backups: "0" or "1" (postgres).
	BackupEnabled string `json:"backup_enabled,omitempty"`
	// ScheduledBackup is a 6-field cron expression
	// (second minute hour day month weekday).
	ScheduledBackup string `json:"scheduled_backup,omitempty"`
	// Instances is 1 for standalone, or 3/5/7 for an HA cluster (postgres).
	Instances *int64 `json:"instances,omitempty"`
}

// RestoreBackupRequest is the body for AddonsService.RestoreBackup and
// ServicesService.RestoreBackup. All fields are optional.
// Schemas: restoreAppAddonBackup, restoreServiceBackup (identical).
type RestoreBackupRequest struct {
	// BackupName selects a specific backup to restore from.
	BackupName string `json:"backup_name,omitempty"`
	// TargetTime is an ISO 8601 timestamp for point-in-time recovery.
	TargetTime string `json:"target_time,omitempty"`
	// ToNewCluster set to "1" restores to a new cluster instead of replacing
	// the current one.
	ToNewCluster string `json:"to_new_cluster,omitempty"`
}

// CreateReplicaRequest is the body for AddonsService.CreateReplica and
// ServicesService.CreateReplica. All fields default to the primary's values.
// Schemas: createAppAddonReplica, createServiceReplica (identical).
type CreateReplicaRequest struct {
	CPUSize *float64 `json:"cpu_size,omitempty"` // cores
	RAMSize *int64   `json:"ram_size,omitempty"` // MB
}

// UpdateAppVarRequest is the body for VarsService.Update.
// Schema: updateAppVar.
type UpdateAppVarRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	// ProjectVariablesEnabled toggles inheritance of project-level
	// environment variables.
	ProjectVariablesEnabled *bool `json:"project_variables_enabled,omitempty"`
}

// CreateAppCronjobRequest is the body for CronjobsService.Create.
// Schema: createAppCronjob.
type CreateAppCronjobRequest struct {
	Label string `json:"label,omitempty"`
	// ScheduleType is cron or interval.
	ScheduleType string `json:"schedule_type,omitempty"`
	// IntervalType is every_10_minutes, hourly or daily (interval schedule).
	IntervalType string `json:"interval_type,omitempty"`
	// Cron is a cron expression, e.g. "0 * * * *" (cron schedule).
	Cron    string `json:"cron,omitempty"`
	Command string `json:"command,omitempty"`
	// DailyTime is the execution time for daily jobs, HH:MM (24-hour).
	DailyTime string `json:"daily_time,omitempty"`
	Minute    string `json:"minute,omitempty"` // 0-59
	Hour      string `json:"hour,omitempty"`   // 0-23
}

// UpdateAppCronjobRequest is the body for CronjobsService.Update.
// Schema: updateAppCronjob.
type UpdateAppCronjobRequest struct {
	Label   string `json:"label,omitempty"`
	Command string `json:"command,omitempty"`
}

// CreateAppPortRequest is the body for PortsService.Create.
// Schema: createAppPort.
type CreateAppPortRequest struct {
	InternalPort int64 `json:"internal_port"` // 1-65535
	// Protocol is tcp or udp.
	Protocol string `json:"protocol"`
	// Public exposes the port publicly.
	Public *bool `json:"public,omitempty"`
}

// DeploymentListOptions are the query parameters for
// DeploymentsService.List.
type DeploymentListOptions struct {
	Status string // filter by deployment status
	Period string // filter by time period
}

// CreateBucketRequest is the body for BucketsService.Create.
// Schema: createBucket.
type CreateBucketRequest struct {
	Label      string `json:"label"`
	ResourceID string `json:"resource_id,omitempty"`
	// MigetID is deprecated; use ResourceID instead.
	MigetID string `json:"miget_id,omitempty"`
	// Visibility is public_access or private_access.
	Visibility string   `json:"visibility,omitempty"`
	DiskSize   *float64 `json:"disk_size,omitempty"` // GiB
	ProjectID  string   `json:"project_id,omitempty"`
}

// UpdateBucketRequest is the body for BucketsService.Update.
// Schema: updateBucket.
type UpdateBucketRequest struct {
	ProjectID string `json:"project_id,omitempty"`
	Label     string `json:"label,omitempty"`
	// Visibility is public_access or private_access.
	Visibility string   `json:"visibility,omitempty"`
	DiskSize   *float64 `json:"disk_size,omitempty"` // GiB
}

// BucketObjectListOptions are the query parameters for
// BucketObjectsService.List.
type BucketObjectListOptions struct {
	Prefix string // folder path prefix
	Limit  int    // page size (zero value omits the parameter)
	Cursor string // pagination cursor from a previous page
}

// GenerateUploadURLRequest is the body for
// BucketObjectsService.GenerateUploadURL.
// Schema: generateUploadUrl.
type GenerateUploadURLRequest struct {
	Key         string `json:"key"`  // object key (file path)
	Size        int64  `json:"size"` // file size in bytes
	ContentType string `json:"content_type,omitempty"`
}

// CreateContainerRegistryCredentialRequest is the body for
// RegistryCredentialsService.Create.
// Schema: createContainerRegistryCredential.
type CreateContainerRegistryCredentialRequest struct {
	Name string `json:"name"`
	// Registry is one of docker_hub, github, gitlab, aws_ecr, azure,
	// digitalocean, quay, generic.
	Registry string `json:"registry"`
	Username string `json:"username"`
	// Token is the registry password or access token.
	Token string `json:"token"`
	// RegistryHostname is required for generic, aws_ecr and azure.
	RegistryHostname string `json:"registry_hostname,omitempty"`
	// SkipValidation skips live credential validation (non-production
	// environments only).
	SkipValidation *bool `json:"skip_validation,omitempty"`
}

// UpdateContainerRegistryCredentialRequest is the body for
// RegistryCredentialsService.Update.
// Schema: updateContainerRegistryCredential.
type UpdateContainerRegistryCredentialRequest struct {
	Name     string `json:"name,omitempty"`
	Registry string `json:"registry,omitempty"`
	Username string `json:"username,omitempty"`
	// Token rotates the stored secret.
	Token            string `json:"token,omitempty"`
	RegistryHostname string `json:"registry_hostname,omitempty"`
	SkipValidation   *bool  `json:"skip_validation,omitempty"`
}

// CreateProjectRequest is the body for ProjectsService.Create.
// Schema: createProject.
type CreateProjectRequest struct {
	// Name must be unique within the workspace.
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateProjectRequest is the body for ProjectsService.Update.
// Schema: updateProject.
type UpdateProjectRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ResourceComponentInput selects a resource component and quantity.
type ResourceComponentInput struct {
	// CodeName is the component code name, e.g. extra_ram, extra_cpu,
	// extra_disk.
	CodeName string `json:"code_name"`
	// Quantity is the number of component units. On update it is the total,
	// not an increment.
	Quantity int64 `json:"quantity"`
}

// CreateResourceRequest is the body for ResourcesService.Create.
// Schema: createResource.
type CreateResourceRequest struct {
	// PlanType is dev (development/hobby) or pro (production).
	PlanType string `json:"plan_type"`
	// PlanCodeName is the specific plan identifier, e.g. free, starter.
	PlanCodeName string `json:"plan_code_name"`
	// RegionCode is the deployment region code, e.g. eu-central-1.
	RegionCode string                   `json:"region_code"`
	Components []ResourceComponentInput `json:"components,omitempty"`
}

// UpdateResourceRequest is the body for ResourcesService.Update.
// Schema: updateResource.
type UpdateResourceRequest struct {
	PlanType     string                   `json:"plan_type,omitempty"`
	PlanCodeName string                   `json:"plan_code_name,omitempty"`
	Components   []ResourceComponentInput `json:"components,omitempty"`
}

// CreateServiceRequest is the body for ServicesService.Create.
// Schema: createService.
type CreateServiceRequest struct {
	// ServiceType is postgres or shared_storage.
	ServiceType string `json:"service_type"`
	ProjectID   string `json:"project_id"`
	ResourceID  string `json:"resource_id,omitempty"`
	// MigetID is deprecated; use ResourceID instead.
	MigetID  string   `json:"miget_id,omitempty"`
	Label    string   `json:"label"`
	RAMSize  *float64 `json:"ram_size,omitempty"`  // MiB
	DiskSize *float64 `json:"disk_size,omitempty"` // GiB
	CPUSize  *float64 `json:"cpu_size,omitempty"`  // cores
	// PostgresVersion selects the PostgreSQL version (postgres).
	PostgresVersion string `json:"postgres_version,omitempty"`
	// PublicAccess enables public internet access: "0" or "1" (postgres).
	PublicAccess string `json:"public_access,omitempty"`
	// EnvironmentVariables injects connection environment variables into the
	// parent application.
	EnvironmentVariables *bool `json:"environment_variables,omitempty"`
	// Instances is 1 for standalone, or 3/5/7 for an HA cluster (postgres).
	Instances *int64 `json:"instances,omitempty"`
	// CreationMode is fresh (default) or external_replica.
	CreationMode string `json:"creation_mode,omitempty"`
	// External replica settings (creation_mode external_replica).
	ExternalHost string `json:"external_host,omitempty"`
	ExternalPort *int64 `json:"external_port,omitempty"`
	// AuthType is password or tls.
	AuthType            string `json:"auth_type,omitempty"`
	ReplicationUsername string `json:"replication_username,omitempty"`
	ReplicationPassword string `json:"replication_password,omitempty"`
	CACrt               string `json:"ca_crt,omitempty"`
	TLSCrt              string `json:"tls_crt,omitempty"`
	TLSKey              string `json:"tls_key,omitempty"`
	// S3Enabled enables S3 WAL archive fallback: "0" or "1".
	S3Enabled    string `json:"s3_enabled,omitempty"`
	S3Endpoint   string `json:"s3_endpoint,omitempty"`
	S3Bucket     string `json:"s3_bucket,omitempty"`
	S3Path       string `json:"s3_path,omitempty"`
	S3AccessKey  string `json:"s3_access_key,omitempty"`
	S3SecretKey  string `json:"s3_secret_key,omitempty"`
	S3ServerName string `json:"s3_server_name,omitempty"`
	// MountPoint is the container mount path, e.g. /data (shared_storage).
	MountPoint string `json:"mount_point,omitempty"`
	// StorageAccess is RWO or RWX (shared_storage).
	StorageAccess string `json:"storage_access,omitempty"`
}

// UpdateServiceRequest is the body for ServicesService.Update.
// Schema: updateService.
type UpdateServiceRequest struct {
	ProjectID string   `json:"project_id,omitempty"`
	Label     string   `json:"label,omitempty"`
	RAMSize   *float64 `json:"ram_size,omitempty"`  // MiB
	DiskSize  *float64 `json:"disk_size,omitempty"` // GiB
	CPUSize   *float64 `json:"cpu_size,omitempty"`  // cores
	// PublicAccess enables public internet access: "0" or "1" (postgres).
	PublicAccess string `json:"public_access,omitempty"`
	// BackupEnabled enables automatic backups: "0" or "1" (postgres).
	BackupEnabled string `json:"backup_enabled,omitempty"`
	// ScheduledBackup is a 6-field cron expression
	// (second minute hour day month weekday).
	ScheduledBackup string `json:"scheduled_backup,omitempty"`
	// Instances is 1 for standalone, or 3/5/7 for an HA cluster (postgres).
	Instances *int64 `json:"instances,omitempty"`
}

// MountAppRequest is the body for ServicesService.MountApp.
// Schema: mountAppToService.
type MountAppRequest struct {
	// AppID is the UUID of the application to mount.
	AppID string `json:"app_id"`
	// MountPoint defaults to the service's mount point.
	MountPoint string `json:"mount_point,omitempty"`
	// Label is a human-readable label for the mounted add-on.
	Label   string `json:"label,omitempty"`
	SubPath string `json:"sub_path,omitempty"`
}

type CreateProjectRestrictionRequest struct {
	UserEmail string `json:"user_email,omitempty"`
	RoleName  string `json:"role_name,omitempty"`
}

type LogEvent struct {
	Event string
	Data  string
}
