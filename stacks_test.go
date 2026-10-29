package lyngbak

import (
	"context"
	"net/http"
	"testing"
)

func TestStackRepositoryURL(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/stacks/st-1", 200,
		`{"uuid":"st-1","name":"shop","label":"Shop","state":"running","project_id":null,"compose_path":".","deployment_config":{"type":"github","repository":"acme/shop","repository_url":"https://github.com/acme/shop","branch":"main","auto_deploy_enabled":true,"credential_id":"c-1"},"total_cpu_size":0.5,"total_ram_size":1024,"total_disk_size":2048,"created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z","services":[{"compose_service_name":"web","compose_service_type":"app","status":"running","item_type":"App","item_uuid":"a-1","item_name":"web"}],"latest_deployment":null}`))
	stack, err := client.Stacks.Get(context.Background(), "st-1")
	if err != nil {
		t.Fatal(err)
	}
	cfg := stack.DeploymentConfig
	if cfg == nil || cfg.Type != "github" || cfg.RepositoryURL == nil || *cfg.RepositoryURL != "https://github.com/acme/shop" || cfg.Repository == nil || *cfg.Repository != "acme/shop" {
		t.Errorf("unexpected config: %+v", cfg)
	}
	if stack.ProjectID != nil || stack.LatestDeployment != nil || len(stack.Services) != 1 || stack.Services[0].ItemName == nil {
		t.Errorf("unexpected stack: %+v", stack)
	}
}

func TestStackPublicGitConfig(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/stacks", 200,
		`[{"uuid":"st-2","deployment_config":{"type":"public_git","repository":null,"repository_url":"https://git.example.com/acme/shop.git","branch":"main","auto_deploy_enabled":null,"credential_id":null}}]`))
	stacks, err := client.Stacks.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	cfg := stacks[0].DeploymentConfig
	if cfg.Type != "public_git" || cfg.Repository != nil || cfg.AutoDeployEnabled != nil || cfg.CredentialID != nil || cfg.RepositoryURL == nil || *cfg.RepositoryURL != "https://git.example.com/acme/shop.git" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestStackGetDeploymentUsesStringID(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/stacks/st-1/deployments/d-1", 200,
		`{"uuid":"d-1","state":"failed","subevent":"build","commit_sha":null,"validation_errors":{"x":["bad"]},"published_at":null,"completed_at":"2026-01-02T03:04:05Z","created_at":"2026-01-02T03:04:05Z"}`))
	d, err := client.Stacks.GetDeployment(context.Background(), "st-1", "d-1")
	if err != nil {
		t.Fatal(err)
	}
	if d.UUID != "d-1" || d.CommitSHA != nil || d.PublishedAt != nil || d.CompletedAt == nil || string(d.ValidationErrors) != `{"x":["bad"]}` {
		t.Errorf("unexpected deployment: %+v", d)
	}
}
