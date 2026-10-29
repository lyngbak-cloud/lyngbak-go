package lyngbak

import (
	"encoding/json"
	"testing"
)

func TestAppDecodesReadBack(t *testing.T) {
	raw := `{
		"uuid": "app-1",
		"health_checks": {"liveness_probe_path": "/healthz", "liveness_probe_enabled": true, "readiness_probe_path": null},
		"scaling_profile": {"replicas": 2, "auto_scaling_enabled": true, "auto_min_replicas": null, "auto_max_replicas": 5, "scaling_start_time": "08:00", "scaling_end_time": "18:00", "cpu_threshold": 70, "memory_threshold": 80},
		"vpc_uuid": null,
		"deployment_config": {"type":"container_registry","image_url":"nginx","tag":"1.27","credential_id":null,"command":["nginx"],"args":[]},
		"ram_size": 128.0,
		"cpu_size": 0.1,
		"private_access": true,
		"basic_auth_enabled": true,
		"basic_auth_username": "admin",
		"project_variables_enabled": true,
		"project_secret_files_enabled": false,
		"internal_url": "app.internal",
		"public_url": null,
		"resource_id": "res-1"
	}`
	var app App
	if err := json.Unmarshal([]byte(raw), &app); err != nil {
		t.Fatal(err)
	}
	if app.HealthChecks == nil || app.HealthChecks.LivenessProbePath == nil || *app.HealthChecks.LivenessProbePath != "/healthz" {
		t.Errorf("liveness path = %+v", app.HealthChecks)
	}
	if app.HealthChecks.ReadinessProbePath != nil {
		t.Errorf("readiness path should be nil")
	}
	if app.HealthChecks.StartupProbePath != nil {
		t.Errorf("absent key should be nil")
	}
	sp := app.ScalingProfile
	if sp == nil || sp.Replicas != 2 || !sp.AutoScalingEnabled || sp.AutoMinReplicas != nil || sp.AutoMaxReplicas == nil || *sp.AutoMaxReplicas != 5 {
		t.Errorf("scaling profile = %+v", sp)
	}
	if *sp.ScalingStartTime != "08:00" || *sp.CPUThreshold != 70 || *sp.MemoryThreshold != 80 {
		t.Errorf("scaling profile = %+v", sp)
	}
	if app.VPCUUID != nil {
		t.Errorf("vpc_uuid should be nil")
	}
	dc := app.DeploymentConfig
	if dc == nil || dc.Type != "container_registry" || dc.CredentialID != nil || dc.Tag == nil || *dc.Tag != "1.27" {
		t.Errorf("deployment config = %+v", dc)
	}
	if len(dc.Command) != 1 || dc.Command[0] != "nginx" || dc.Args == nil || len(dc.Args) != 0 {
		t.Errorf("command/args = %v / %v", dc.Command, dc.Args)
	}
	if app.RAMSize != 128 || app.CPUSize != 0.1 {
		t.Errorf("sizes = %v / %v", app.RAMSize, app.CPUSize)
	}
	if !app.PrivateAccess || !app.BasicAuthEnabled || *app.BasicAuthUsername != "admin" || !app.ProjectVariablesEnabled || app.ProjectSecretFilesEnabled {
		t.Errorf("flags = %+v", app)
	}
	if app.InternalURL == nil || app.PublicURL != nil || app.ResourceID != "res-1" {
		t.Errorf("urls = %v / %v", app.InternalURL, app.PublicURL)
	}
}

func TestAppNullDeploymentConfig(t *testing.T) {
	var app App
	if err := json.Unmarshal([]byte(`{"deployment_config": null}`), &app); err != nil {
		t.Fatal(err)
	}
	if app.DeploymentConfig != nil {
		t.Errorf("deployment config should be nil")
	}
}

func TestAddonServiceIDIsInteger(t *testing.T) {
	var a Addon
	if err := json.Unmarshal([]byte(`{"service_id": 42}`), &a); err != nil {
		t.Fatal(err)
	}
	if a.ServiceID == nil || *a.ServiceID != 42 {
		t.Errorf("service id = %v", a.ServiceID)
	}
	var b Addon
	if err := json.Unmarshal([]byte(`{"service_id": null}`), &b); err != nil {
		t.Fatal(err)
	}
	if b.ServiceID != nil {
		t.Errorf("service id should be nil")
	}
}

func TestAddonDecodesReplicasAndBackup(t *testing.T) {
	raw := `{"ram_size":256,"disk_size":1.5,"cpu_size":0.25,"backup_enabled":null,"scheduled_backup":null,"role":"primary","primary_addon_uuid":null,"replicas":[{"uuid":"r1","name":"n","label":"l","state":"running"}]}`
	var a Addon
	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		t.Fatal(err)
	}
	if a.RAMSize != 256 || a.DiskSize != 1.5 || a.CPUSize != 0.25 || a.BackupEnabled != nil || a.ScheduledBackup != nil {
		t.Errorf("addon = %+v", a)
	}
	if a.Role == nil || *a.Role != "primary" || a.PrimaryAddonUUID != nil || len(a.Replicas) != 1 || a.Replicas[0].UUID != "r1" {
		t.Errorf("addon = %+v", a)
	}
}

func TestServiceDecodesMounts(t *testing.T) {
	raw := `{"mounts":[{"app_id":"app-1","addon_uuid":"ad-1","mount_point":"/data","sub_path":null,"label":null}],"miget_id":null,"instances":null,"backup_enabled":null,"public_access":false,"postgres_version":"16"}`
	var s Service
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Mounts) != 1 {
		t.Fatalf("mounts = %+v", s.Mounts)
	}
	m := s.Mounts[0]
	if m.SubPath != nil || m.Label != nil || m.AppID != "app-1" || m.MountPoint != "/data" || m.AddonUUID != "ad-1" {
		t.Errorf("mount = %+v", m)
	}
	if s.MigetID != nil || s.Instances != nil || s.BackupEnabled != nil || s.PublicAccess == nil || *s.PublicAccess || *s.PostgresVersion != "16" {
		t.Errorf("service = %+v", s)
	}
}

func TestPortFlatEndpoint(t *testing.T) {
	raw := `{"id":42,"internal_port":8080,"external_port":null,"public":true,"endpoint_id":7,"endpoint_hostname":"p.example.net","endpoint_ipv4_primary":"10.0.0.1"}`
	var p Port
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	if p.EndpointHostname != "p.example.net" || p.EndpointIPv4Primary != "10.0.0.1" || p.EndpointID != 7 {
		t.Errorf("port = %+v", p)
	}
	if p.ID != 42 || p.InternalPort != 8080 || p.ExternalPort != nil || !p.Public {
		t.Errorf("port = %+v", p)
	}
}

func TestCronjobCreateHasNoName(t *testing.T) {
	b, err := json.Marshal(CreateAppCronjobRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["name"]; ok {
		t.Errorf("name key present: %s", b)
	}
}

func TestResourceProjectBucketDecode(t *testing.T) {
	var r Resource
	if err := json.Unmarshal([]byte(`{"region_code":"eu-central-1","assigned":true,"project_ids":["p1"],"component_quantities":[{"code_name":"extra_ram","quantity":2}],"plan":{"plan_type":"pro"}}`), &r); err != nil {
		t.Fatal(err)
	}
	if r.RegionCode != "eu-central-1" || !r.Assigned || len(r.ProjectIDs) != 1 || r.ComponentQuantities[0].CodeName != "extra_ram" || r.ComponentQuantities[0].Quantity != 2 || r.Plan.PlanType != "pro" {
		t.Errorf("resource = %+v", r)
	}
	var p Project
	if err := json.Unmarshal([]byte(`{"uuid":"0b0f2f9e-0000-4000-8000-000000000000","resources":[{"uuid":"r1"}],"restrictions":[{"id":3,"subject_type":"user","user_uuid":"u1","role_name":null}]}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.UUID == "" || len(p.Resources) != 1 || p.Restrictions[0].ID != 3 || p.Restrictions[0].UserUUID == nil || p.Restrictions[0].RoleName != nil {
		t.Errorf("project = %+v", p)
	}
	var b Bucket
	if err := json.Unmarshal([]byte(`{"resource_id":"r1","project_id":null}`), &b); err != nil {
		t.Fatal(err)
	}
	if b.ResourceID != "r1" || b.ProjectID != nil {
		t.Errorf("bucket = %+v", b)
	}
}

func TestDeploymentCronjobBackupDecode(t *testing.T) {
	var d Deployment
	if err := json.Unmarshal([]byte(`{"branch":"main","commit_sha":null,"commit_message":null}`), &d); err != nil {
		t.Fatal(err)
	}
	if d.Branch == nil || *d.Branch != "main" || d.CommitSHA != nil || d.CommitMessage != nil {
		t.Errorf("deployment = %+v", d)
	}
	var c Cronjob
	if err := json.Unmarshal([]byte(`{"label":"L","cron_expression":"0 * * * *"}`), &c); err != nil {
		t.Fatal(err)
	}
	if c.Label != "L" || c.CronExpression != "0 * * * *" {
		t.Errorf("cronjob = %+v", c)
	}
	var bk Backup
	if err := json.Unmarshal([]byte(`{"backup_active":true}`), &bk); err != nil {
		t.Fatal(err)
	}
	if !bk.BackupActive {
		t.Errorf("backup = %+v", bk)
	}
}

func TestAppCronjobsAreAddonShaped(t *testing.T) {
	raw := `{"uuid":"app-1","cronjobs":[{"uuid":"cj-1","name":"cron-1","type":"cronjob","label":"live","state":"active","quota":{"cpu_size":0.1},"settings":{"schedule_type":"interval","interval_type":"hourly","minute":"0","hour":"","command":"echo live","cron_expression":"0 * * * *"},"var":null,"connection_details":{},"ram_size":64,"disk_size":0,"cpu_size":0.1,"service_id":null,"created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}]}`
	var app App
	if err := json.Unmarshal([]byte(raw), &app); err != nil {
		t.Fatal(err)
	}
	if len(app.Cronjobs) != 1 || app.Cronjobs[0].UUID != "cj-1" || app.Cronjobs[0].Type != "cronjob" || app.Cronjobs[0].Settings["command"] != "echo live" || app.Cronjobs[0].CPUSize != 0.1 {
		t.Errorf("unexpected cronjobs: %+v", app.Cronjobs)
	}
}

func TestServiceDecodesRoleReplicasAndAddons(t *testing.T) {
	raw := `{"uuid":"svc-1","role":"primary","primary_addon_uuid":null,"replicas":[{"uuid":"svc-2","role":"replica","primary_addon_uuid":"ad-1"}],"addons":[{"uuid":"ad-1","name":"pg-1","label":"db","state":"active","is_replica":false,"settings":{"postgres_version":"16"},"created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}]}`
	var svc Service
	if err := json.Unmarshal([]byte(raw), &svc); err != nil {
		t.Fatal(err)
	}
	if svc.Role == nil || *svc.Role != "primary" || svc.PrimaryAddonUUID != nil || len(svc.Replicas) != 1 || svc.Replicas[0].PrimaryAddonUUID == nil || *svc.Replicas[0].PrimaryAddonUUID != "ad-1" {
		t.Errorf("unexpected service: %+v", svc)
	}
	if len(svc.Addons) != 1 || svc.Addons[0].UUID != "ad-1" || svc.Addons[0].Settings["postgres_version"] != "16" || svc.Addons[0].IsReplica {
		t.Errorf("unexpected addons: %+v", svc.Addons)
	}
}

func marshalKeys(t *testing.T, v any) map[string]json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func requireKeys(t *testing.T, got map[string]json.RawMessage, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("keys = %d, want %d: %v", len(got), len(want), got)
	}
	for _, k := range want {
		if _, ok := got[k]; !ok {
			t.Errorf("missing key %q", k)
		}
	}
}

func TestAppDeploymentConfigKeys(t *testing.T) {
	s, b := "x", true
	got := marshalKeys(t, AppDeploymentConfig{
		RepositoryURL: s, ProjectPath: s, RunCommand: s, Language: s, BuildCommand: s,
		PreDeployCommand: s, PostDeployCommand: s, UseDHI: &b,
	})
	requireKeys(t, got, "repository_url", "project_path", "run_command", "language", "build_command", "pre_deploy_command", "post_deploy_command", "use_dhi")
}

func TestUpdateAppDeploymentConfigKeys(t *testing.T) {
	s, b := "x", true
	got := marshalKeys(t, UpdateAppDeploymentConfig{
		RepositoryURL: s, Repository: s, Branch: &s, AutoDeployEnabled: &b, AutoDeployBranch: &s,
		DockerfilePath: &s, BuildContext: &s, ProjectPath: &s, RunCommand: &s, Language: &s,
		BuildCommand: &s, PreDeployCommand: &s, PostDeployCommand: &s, UseDHI: &b,
	})
	requireKeys(t, got, "repository_url", "repository", "branch", "auto_deploy_enabled", "auto_deploy_branch", "dockerfile_path", "build_context", "project_path", "run_command", "language", "build_command", "pre_deploy_command", "post_deploy_command", "use_dhi")
}

func TestClearableStringsMarshalEmpty(t *testing.T) {
	empty := Ptr("")
	cases := map[string]struct {
		v   any
		key string
	}{
		"project description": {UpdateProjectRequest{Description: empty}, "description"},
		"scaling start":       {UpdateAppScalingProfileRequest{ScalingStartTime: empty}, "scaling_start_time"},
		"scaling end":         {UpdateAppScalingProfileRequest{ScalingEndTime: empty}, "scaling_end_time"},
		"trigger label":       {UpdatePreviewEnvironmentConfigRequest{TriggerLabel: empty}, "trigger_label"},
		"deployment branch":   {UpdateAppDeploymentConfig{Branch: empty}, "branch"},
	}
	for name, c := range cases {
		got := marshalKeys(t, c.v)
		if string(got[c.key]) != `""` {
			t.Errorf("%s: %s = %s, want \"\"", name, c.key, got[c.key])
		}
	}
	requireKeys(t, marshalKeys(t, UpdateProjectRequest{}))
}

func TestAssignProjectResourceRequestKeys(t *testing.T) {
	requireKeys(t, marshalKeys(t, AssignProjectResourceRequest{ResourceID: "r"}), "resource_id")
	requireKeys(t, marshalKeys(t, AssignProjectResourceRequest{ResourceID: "r", WithHostingProjects: Ptr(false)}), "resource_id", "with_hosting_projects")
}
