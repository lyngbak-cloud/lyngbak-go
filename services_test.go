package lyngbak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// route asserts the request method and path, then writes the given response
// body with the given status.
func route(t *testing.T, method, path string, status int, body string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != path {
			t.Errorf("path = %s, want %s", r.URL.Path, path)
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}
}

func TestAppsCreate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/apps" {
			t.Errorf("got %s %s, want POST /api/v1/apps", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["name"] != "web" || body["resource_id"] != "res-1" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"uuid":"app-1","name":"web","label":"Web","state":"pending","project_id":"proj-1"}`)
	})
	app, err := client.Apps.Create(context.Background(), &CreateAppRequest{
		Label: "Web", Name: "web", ProjectID: "proj-1", ResourceID: "res-1", Builder: "auto",
	})
	if err != nil {
		t.Fatal(err)
	}
	if app.UUID != "app-1" || app.State != "pending" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestAppsUpdateState(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/apps/app-1/state" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["state"] != AppStateScheduleRestart {
			t.Errorf("state = %q, want %q", body["state"], AppStateScheduleRestart)
		}
		fmt.Fprint(w, `{"message":"App restart scheduled"}`)
	})
	msg, err := client.Apps.UpdateState(context.Background(), "app-1", AppStateScheduleRestart)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Message != "App restart scheduled" {
		t.Errorf("message = %q", msg.Message)
	}
}

func TestProjectsList(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/projects", 200,
		`[{"uuid":"proj-7","name":"demo","description":"d"}]`))
	projects, err := client.Projects.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 1 || projects[0].UUID != "proj-7" || projects[0].Name != "demo" {
		t.Errorf("unexpected projects: %+v", projects)
	}
}

func TestProjectVarsUpdate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/projects/proj-1/vars" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "K" || body["value"] != "V" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusOK)
	})
	if err := client.ProjectVars.Update(context.Background(), "proj-1", "K", "V"); err != nil {
		t.Fatal(err)
	}
}

func TestResourcesGet(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/resources/res-1", 200,
		`{"uuid":"res-1","name":"miget-1","state":"active","labels":["prod"],"total_ram_size":2048,"plan":{"code_name":"starter"}}`))
	res, err := client.Resources.Get(context.Background(), "res-1")
	if err != nil {
		t.Fatal(err)
	}
	if res.UUID != "res-1" || res.TotalRAMSize != 2048 || res.Plan.CodeName != "starter" {
		t.Errorf("unexpected resource: %+v", res)
	}
}

func TestServicesMountApp(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/services/svc-1/mount_app" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["app_id"] != "app-1" || body["mount_point"] != "/data" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"uuid":"svc-1","service_type":"shared_storage"}`)
	})
	svc, err := client.Services.MountApp(context.Background(), "svc-1", &MountAppRequest{AppID: "app-1", MountPoint: "/data"})
	if err != nil {
		t.Fatal(err)
	}
	if svc.ServiceType != "shared_storage" {
		t.Errorf("unexpected service: %+v", svc)
	}
}

func TestBucketsUpdatePolicy(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/buckets/b-1/policy" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["policy"] != `{"Version":"2012-10-17"}` {
			t.Errorf("policy = %q", body["policy"])
		}
		fmt.Fprint(w, `{"uuid":"b-1","policy":"{\"Version\":\"2012-10-17\"}"}`)
	})
	b, err := client.Buckets.UpdatePolicy(context.Background(), "b-1", `{"Version":"2012-10-17"}`)
	if err != nil {
		t.Fatal(err)
	}
	if b.UUID != "b-1" {
		t.Errorf("unexpected bucket: %+v", b)
	}
}

func TestBucketObjectsList(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/buckets/b-1/objects/list" {
			t.Errorf("path = %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("prefix") != "img/" || q.Get("limit") != "50" || q.Get("cursor") != "c1" {
			t.Errorf("unexpected query: %v", q)
		}
		fmt.Fprint(w, `{"objects":[{"key":"img/a.png","name":"a.png","is_folder":false,"size":12,"content_type":"image/png"}],"current_path":"img/","has_more":false}`)
	})
	list, err := client.BucketObjects.List(context.Background(), "b-1", &BucketObjectListOptions{Prefix: "img/", Limit: 50, Cursor: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Objects) != 1 || list.Objects[0].Key != "img/a.png" || list.HasMore {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestBucketObjectsKeyEscaping(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// The raw path must contain the escaped key as a single segment.
		want := "/api/v1/buckets/b-1/objects/" + "img%2Fmy%20file.png"
		if r.URL.EscapedPath() != want {
			t.Errorf("escaped path = %s, want %s", r.URL.EscapedPath(), want)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.BucketObjects.Delete(context.Background(), "b-1", "img/my file.png"); err != nil {
		t.Fatal(err)
	}
}

func TestAddonsCreate(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPost, "/api/v1/apps/app-1/addons", 201,
		`{"uuid":"ad-1","type":"postgres","state":"pending","var":{"DATABASE_URL":"postgres://x"}}`))
	addon, err := client.Addons.Create(context.Background(), "app-1", &CreateAppAddonRequest{
		Type: "postgres", PostgresVersion: "16", Instances: Ptr(int64(3)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if addon.UUID != "ad-1" || addon.Var["DATABASE_URL"] != "postgres://x" {
		t.Errorf("unexpected addon: %+v", addon)
	}
}

func TestDeploymentsList(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/apps/app-1/deployments" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "completed" {
			t.Errorf("status query = %q", r.URL.Query().Get("status"))
		}
		fmt.Fprint(w, `[{"uuid":"d-1","state":"completed"}]`)
	})
	deps, err := client.Deployments.List(context.Background(), "app-1", &DeploymentListOptions{Status: "completed"})
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 1 || deps[0].State != "completed" {
		t.Errorf("unexpected deployments: %+v", deps)
	}
}

func TestDomainsVerify(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPost, "/api/v1/apps/app-1/domains/dom-1/verify", 201,
		`{"uuid":"dom-1","name":"app.example.com","verification_status":"in_progress"}`))
	d, err := client.Domains.Verify(context.Background(), "app-1", "dom-1")
	if err != nil {
		t.Fatal(err)
	}
	if d.VerificationStatus != "in_progress" {
		t.Errorf("unexpected domain: %+v", d)
	}
}

func TestDomainsCreateNestedBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["domain"]["name"] != "app.example.com" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"uuid":"dom-1","name":"app.example.com"}`)
	})
	if _, err := client.Domains.Create(context.Background(), "app-1", "app.example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestVarsDelete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/apps/app-1/vars" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("key") != "API_KEY" {
			t.Errorf("key query = %q", r.URL.Query().Get("key"))
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.Vars.Delete(context.Background(), "app-1", "API_KEY"); err != nil {
		t.Fatal(err)
	}
}

func TestCronjobsGet(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/apps/app-1/cronjobs/cj-1", 200,
		`{"uuid":"cj-1","name":"cleanup","schedule_type":"interval","interval_type":"daily","command":"rake cleanup"}`))
	cj, err := client.Cronjobs.Get(context.Background(), "app-1", "cj-1")
	if err != nil {
		t.Fatal(err)
	}
	if cj.Command != "rake cleanup" {
		t.Errorf("unexpected cronjob: %+v", cj)
	}
}

func TestPortsExposePublicly(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPatch, "/api/v1/apps/app-1/ports/42/expose_publicly", 200,
		`{"id":42,"internal_port":8080,"protocol":"tcp","public":true,"endpoint_hostname":"p.miget.net"}`))
	p, err := client.Ports.ExposePublicly(context.Background(), "app-1", 42)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Public || p.EndpointHostname != "p.miget.net" {
		t.Errorf("unexpected port: %+v", p)
	}
}

func TestAuthSignIn(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/sign_in" {
			t.Errorf("path = %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["email"] != "a@b.c" || body["password"] != "pw" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"access_token":"jwt-a","refresh_token":"jwt-r"}`)
	})
	tokens, err := client.Auth.SignIn(context.Background(), "a@b.c", "pw")
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken != "jwt-a" || tokens.RefreshToken != "jwt-r" {
		t.Errorf("unexpected tokens: %+v", tokens)
	}
}

func TestUsersMe(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/users/me", 200,
		`{"uuid":"u-1","email":"a@b.c","name":"Ada"}`))
	u, err := client.Users.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if u.Email != "a@b.c" {
		t.Errorf("unexpected user: %+v", u)
	}
}

func TestSSHKeysCreateAndDelete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/users/me/ssh_keys":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":5,"public_key":"ssh-ed25519 AAAA"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/users/me/ssh_keys/5":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})
	key, err := client.SSHKeys.Create(context.Background(), "ssh-ed25519 AAAA")
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 5 {
		t.Errorf("unexpected key: %+v", key)
	}
	if err := client.SSHKeys.Delete(context.Background(), key.ID); err != nil {
		t.Fatal(err)
	}
}

func TestCreditsGet(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/users/me/credits", 200,
		`{"credit_balance_in_cents":1500,"credit_balance_formatted":"$15.00","referral_code":"ref"}`))
	c, err := client.Credits.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c.CreditBalanceInCents != 1500 {
		t.Errorf("unexpected credits: %+v", c)
	}
}

func TestPlansList(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/plans", 200,
		`[{"name":"Starter","code_name":"starter","ram_size":512,"unit_price":500}]`))
	plans, err := client.Plans.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].CodeName != "starter" {
		t.Errorf("unexpected plans: %+v", plans)
	}
}

func TestRegionsList(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/regions", 200,
		`[{"id":1,"name":"EU Central","code":"eu-central-1"}]`))
	regions, err := client.Regions.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(regions) != 1 || regions[0].Code != "eu-central-1" {
		t.Errorf("unexpected regions: %+v", regions)
	}
}

func TestComponentsList(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/components", 200,
		`[{"name":"Extra RAM","code_name":"extra_ram","unit_price":100}]`))
	comps, err := client.Components.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(comps) != 1 || comps[0].CodeName != "extra_ram" {
		t.Errorf("unexpected components: %+v", comps)
	}
}

func TestRegistryCredentialsCreate(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPost, "/api/v1/container_registry_credentials", 201,
		`{"uuid":"cred-1","name":"docker","registry":"docker_hub","username":"u"}`))
	cred, err := client.RegistryCredentials.Create(context.Background(), &CreateContainerRegistryCredentialRequest{
		Name: "docker", Registry: "docker_hub", Username: "u", Token: "t",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cred.UUID != "cred-1" || cred.Registry != "docker_hub" {
		t.Errorf("unexpected credential: %+v", cred)
	}
}
