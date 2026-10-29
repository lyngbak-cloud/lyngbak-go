package lyngbak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestClient returns a Client pointed at an httptest server running the
// given handler. The server is shut down when the test finishes.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	opts = append([]Option{WithBaseURL(srv.URL)}, opts...)
	return New("test-token", opts...)
}

func TestAuthorizationHeader(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer test-token")
		}
		fmt.Fprint(w, `[]`)
	})
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceHeaderSet(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Workspace-Id"); got != "ws-123" {
			t.Errorf("X-Workspace-Id header = %q, want %q", got, "ws-123")
		}
		fmt.Fprint(w, `[]`)
	}, WithWorkspaceID("ws-123"))
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceHeaderAbsentByDefault(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["X-Workspace-Id"]; ok {
			t.Error("X-Workspace-Id header should not be set")
		}
		fmt.Fprint(w, `[]`)
	})
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestUserAgent(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		want := "lyngbak-go/" + Version
		if got := r.Header.Get("User-Agent"); got != want {
			t.Errorf("User-Agent = %q, want %q", got, want)
		}
		fmt.Fprint(w, `[]`)
	})
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWithUserAgent(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "custom-agent/1.0" {
			t.Errorf("User-Agent = %q, want %q", got, "custom-agent/1.0")
		}
		fmt.Fprint(w, `[]`)
	}, WithUserAgent("custom-agent/1.0"))
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestBaseURLAndPath(t *testing.T) {
	var gotPath string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		fmt.Fprint(w, `{"uuid":"app-1"}`)
	})
	if _, err := client.Apps.Get(context.Background(), "app-1"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v1/apps/app-1" {
		t.Errorf("request path = %q, want %q", gotPath, "/api/v1/apps/app-1")
	}
}

func TestBaseURLTrailingSlashTrimmed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "//") {
			t.Errorf("request path %q contains a double slash", r.URL.Path)
		}
		fmt.Fprint(w, `[]`)
	}))
	t.Cleanup(srv.Close)
	client := New("tok", WithBaseURL(srv.URL+"/"))
	if _, err := client.Apps.List(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorSingleErrorShape(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"App not found"}`)
	})
	_, err := client.Apps.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %T is not *APIError", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Message != "App not found" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "App not found")
	}
	if !strings.Contains(apiErr.Error(), "App not found") {
		t.Errorf("Error() = %q, missing message", apiErr.Error())
	}
	if !IsNotFound(err) {
		t.Error("IsNotFound = false, want true")
	}
}

func TestAPIErrorErrorsArrayShape(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"errors":["name is taken","label is too long"]}`)
	})
	_, err := client.Projects.Create(context.Background(), &CreateProjectRequest{Name: "x"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %T is not *APIError", err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("StatusCode = %d, want 422", apiErr.StatusCode)
	}
	want := []string{"name is taken", "label is too long"}
	if len(apiErr.Errors) != 2 || apiErr.Errors[0] != want[0] || apiErr.Errors[1] != want[1] {
		t.Errorf("Errors = %v, want %v", apiErr.Errors, want)
	}
	if !strings.Contains(apiErr.Error(), "name is taken; label is too long") {
		t.Errorf("Error() = %q, missing joined messages", apiErr.Error())
	}
	if IsNotFound(err) {
		t.Error("IsNotFound = true for 422, want false")
	}
}

func TestAPIErrorRawBodyFallback(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `upstream timeout`)
	})
	_, err := client.Apps.List(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %T is not *APIError", err)
	}
	if string(apiErr.Body) != "upstream timeout" {
		t.Errorf("Body = %q, want raw body preserved", apiErr.Body)
	}
	if !strings.Contains(apiErr.Error(), "upstream timeout") {
		t.Errorf("Error() = %q, missing raw body fallback", apiErr.Error())
	}
}

func TestAPIErrorEmptyBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	_, err := client.Apps.List(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %T is not *APIError", err)
	}
	if !strings.Contains(apiErr.Error(), "401") || !strings.Contains(apiErr.Error(), "Unauthorized") {
		t.Errorf("Error() = %q, want status text fallback", apiErr.Error())
	}
}

func TestNewFromEnv(t *testing.T) {
	t.Setenv(EnvAPIToken, "env-token")
	client, err := NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if client.token != "env-token" {
		t.Errorf("token = %q, want %q", client.token, "env-token")
	}
}

func TestNewFromEnvMissing(t *testing.T) {
	t.Setenv(EnvAPIToken, "")
	if _, err := NewFromEnv(); err == nil {
		t.Fatal("expected error when MIGET_API_TOKEN is unset")
	}
}

func TestRequestBodyEncoding(t *testing.T) {
	var got map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"uuid":"proj-12","name":"demo"}`)
	})
	_, err := client.Projects.Create(context.Background(), &CreateProjectRequest{Name: "demo"})
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "demo" {
		t.Errorf("request body name = %v, want demo", got["name"])
	}
	if _, ok := got["description"]; ok {
		t.Error("empty optional field description should be omitted")
	}
}

func TestPtr(t *testing.T) {
	if v := Ptr(true); *v != true {
		t.Error("Ptr(true) mismatch")
	}
	if v := Ptr(2.5); *v != 2.5 {
		t.Error("Ptr(2.5) mismatch")
	}
}

func TestWithBaseURLStripsAPIPrefix(t *testing.T) {
	for _, in := range []string{"https://api.example.test/api/v1", "https://api.example.test/api/v1/", "https://api.example.test"} {
		if got := New("t", WithBaseURL(in)).baseURL; got != "https://api.example.test" {
			t.Errorf("WithBaseURL(%q) base = %q", in, got)
		}
	}
}

func TestNewFromEnvReadsWorkspaceAndURL(t *testing.T) {
	t.Setenv("MIGET_API_TOKEN", "tok")
	t.Setenv("MIGET_WORKSPACE_ID", "ws-1")
	t.Setenv("MIGET_API_URL", "https://api.example.test/api/v1")
	c, err := NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if c.workspaceID != "ws-1" || c.baseURL != "https://api.example.test" {
		t.Errorf("workspace %q base %q", c.workspaceID, c.baseURL)
	}
	c, _ = NewFromEnv(WithWorkspaceID("ws-2"))
	if c.workspaceID != "ws-2" {
		t.Errorf("explicit option lost: %q", c.workspaceID)
	}
}
