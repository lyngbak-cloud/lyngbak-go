package lyngbak

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSecretFileNameIsEscaped(t *testing.T) {
	cases := []struct{ name, want string }{
		{"my file.txt", "/api/v1/apps/app-1/secret_files/my%20file.txt"},
		{".env.production", "/api/v1/apps/app-1/secret_files/.env.production"},
	}
	for _, tc := range cases {
		var got string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			got = r.URL.EscapedPath()
			w.Write([]byte(`{"name":"sf-1","filename":"/app/x"}`))
		})
		if _, err := client.SecretFiles.Get(context.Background(), "app-1", tc.name); err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("path = %s, want %s", got, tc.want)
		}
	}
}

func TestProjectSecretFileCreate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/projects/p-1/secret_files" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body) != 2 || body["filename"] != "/app/config/example.json" || body["text"] != "placeholder" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(201)
		w.Write([]byte(`{"name":"sf-1","filename":"/app/config/example.json","created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}`))
	})
	sf, err := client.ProjectSecretFiles.Create(context.Background(), "p-1", &SecretFileRequest{Filename: "/app/config/example.json", Text: "placeholder"})
	if err != nil {
		t.Fatal(err)
	}
	if sf.Name != "sf-1" || sf.Filename != "/app/config/example.json" || sf.CreatedAt.IsZero() {
		t.Errorf("unexpected secret file: %+v", sf)
	}
}

func TestGitCredentialsGet(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/git_credentials/gc-1", 200,
		`{"uuid":"gc-1","name":"gh","provider":"github","installation_id":"42","created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z"}`))
	gc, err := client.GitCredentials.Get(context.Background(), "gc-1")
	if err != nil {
		t.Fatal(err)
	}
	if gc.UUID != "gc-1" || gc.Provider != "github" || gc.InstallationID != "42" {
		t.Errorf("unexpected credential: %+v", gc)
	}
}
