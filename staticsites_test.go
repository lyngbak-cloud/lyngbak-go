package lyngbak

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStaticSitesListIsArray(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/static", 200, `[{"uuid":"s-1"},{"uuid":"s-2"}]`))
	sites, err := client.StaticSites.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sites) != 2 || sites[0].UUID != "s-1" || sites[1].UUID != "s-2" {
		t.Errorf("unexpected sites: %+v", sites)
	}
}

func TestStaticSiteDecodesDeploymentConfig(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/static/s-1", 200,
		`{"uuid":"s-1","name":"docs","label":"Docs","state":"running","url":"https://docs.example.com","region":"eu-east-1","project_id":null,"deployment_config":{"source_type":"github","credential_id":"c-1","spa_mode":true,"auto_deploy_enabled":false,"branch":"main","build_command":null},"created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}`))
	site, err := client.StaticSites.Get(context.Background(), "s-1")
	if err != nil {
		t.Fatal(err)
	}
	cfg := site.DeploymentConfig
	if cfg == nil || cfg.SourceType != "github" || cfg.CredentialID == nil || *cfg.CredentialID != "c-1" || !cfg.SPAMode || cfg.BuildCommand != nil {
		t.Errorf("unexpected deployment config: %+v", cfg)
	}
	if site.ProjectID != nil || site.CreatedAt.IsZero() {
		t.Errorf("unexpected site: %+v", site)
	}
}

func TestStaticSiteDomainRoutes(t *testing.T) {
	var gotMethod, gotPath, gotBody string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Write([]byte(`{"uuid":"d-1","name":"www.example.com"}`))
	})
	d, err := client.StaticSites.CreateDomain(context.Background(), "s-1", "www.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost || gotPath != "/api/v1/static/s-1/domains" {
		t.Errorf("got %s %s", gotMethod, gotPath)
	}
	if gotBody != `{"domain":{"name":"www.example.com"}}` || d.UUID != "d-1" {
		t.Errorf("body = %s, domain = %+v", gotBody, d)
	}
}

func TestStaticSiteDeployWithoutArchive(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/static/s-1/deployments" || len(b) != 0 {
			t.Errorf("got %s %s body %q", r.Method, r.URL.Path, b)
		}
		w.Write([]byte(`{"message":"Deployment started"}`))
	})
	msg, err := client.StaticSites.Deploy(context.Background(), "s-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Message != "Deployment started" {
		t.Errorf("message = %q", msg.Message)
	}
}

func TestStaticSiteDeployWithArchive(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		f, hdr, err := r.FormFile("archive")
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(f)
		if string(b) != "zipbytes" || hdr.Filename != "site.zip" {
			t.Errorf("archive = %q (%s)", b, hdr.Filename)
		}
		w.Write([]byte(`{"message":"Deployment started"}`))
	})
	_, err := client.StaticSites.Deploy(context.Background(), "s-1", &DeployStaticSiteRequest{Archive: strings.NewReader("zipbytes"), Filename: "site.zip"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStaticSiteListFilesQuery(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("prefix") != "assets/" || q.Get("limit") != "50" || q.Get("cursor") != "c1" {
			t.Errorf("query = %v", q)
		}
		w.Write([]byte(`{"objects":[],"current_path":"assets/","has_more":false}`))
	})
	list, err := client.StaticSites.ListFiles(context.Background(), "s-1", &BucketObjectListOptions{Prefix: "assets/", Limit: 50, Cursor: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	if list.CurrentPath != "assets/" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestPreviewEnvironmentIDIsInteger(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/apps/a-1/preview_environments/7", 200, `{"id":7,"status":"active","pull_request_number":null}`))
	env, err := client.PreviewEnvironments.Get(context.Background(), "a-1", 7)
	if err != nil {
		t.Fatal(err)
	}
	if env.ID != 7 || env.Status != "active" || env.PullRequestNumber != nil {
		t.Errorf("unexpected environment: %+v", env)
	}
}

func TestPreviewEnvironmentListStatusQuery(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "failed" {
			t.Errorf("query = %v", r.URL.Query())
		}
		w.Write([]byte(`[{"id":1},{"id":2}]`))
	})
	envs, err := client.PreviewEnvironments.List(context.Background(), "a-1", "failed")
	if err != nil {
		t.Fatal(err)
	}
	if len(envs) != 2 {
		t.Errorf("envs = %+v", envs)
	}
}

func TestPreviewEnvironmentUpdateConfigBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/apps/a-1/preview_environments/config" || string(b) != `{"enabled":false,"retention_days":3}` {
			t.Errorf("got %s %s %s", r.Method, r.URL.Path, b)
		}
		w.Write([]byte(`{"enabled":false,"resource_id":null,"clone_settings":{"clone_variables":true},"retention_days":3}`))
	})
	cfg, err := client.PreviewEnvironments.UpdateConfig(context.Background(), "a-1", &UpdatePreviewEnvironmentConfigRequest{Enabled: Ptr(false), RetentionDays: Ptr(int64(3))})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Enabled || cfg.ResourceID != nil || cfg.RetentionDays != 3 || cfg.CloneSettings["clone_variables"] != true {
		t.Errorf("unexpected config: %+v", cfg)
	}
}
