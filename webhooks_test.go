package lyngbak

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestWebhooksListIsArray(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/webhooks", 200,
		`[{"uuid":"w-1","event_filter":["deploy_ended"],"app_uuids":[]},{"uuid":"w-2"}]`))
	hooks, err := client.Webhooks.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 2 || hooks[0].UUID != "w-1" || hooks[0].EventFilter[0] != "deploy_ended" {
		t.Errorf("unexpected webhooks: %+v", hooks)
	}
}

func TestWebhookDeliveriesListIsArray(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/webhooks/w-1/deliveries", 200,
		`[{"uuid":"dl-1","event_id":"e-1","event_type":"deploy_ended","status":"failed","response_code":500,"message":"boom","attempts":2,"last_attempt_at":"2026-01-02T03:04:05Z","created_at":"2026-01-02T03:04:05Z","payload":{"id":"e-1"}},{"uuid":"dl-2","status":"pending","response_code":null,"message":null,"last_attempt_at":null}]`))
	list, err := client.Webhooks.ListDeliveries(context.Background(), "w-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ResponseCode == nil || *list[0].ResponseCode != 500 || string(list[0].Payload) != `{"id":"e-1"}` {
		t.Errorf("unexpected deliveries: %+v", list)
	}
	if list[1].ResponseCode != nil || list[1].Message != nil || list[1].LastAttemptAt != nil {
		t.Errorf("unexpected nullable fields: %+v", list[1])
	}
}

func TestWebhookCreatedCarriesSecret(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPost, "/api/v1/webhooks", 201,
		`{"uuid":"w-1","name":"ci","url":"https://hooks.example.com/in","enabled":true,"event_filter":["deploy_started"],"app_uuids":[],"include_review_apps":false,"secret":"whsec_example"}`))
	created, err := client.Webhooks.Create(context.Background(), &CreateWebhookRequest{Name: "ci", URL: "https://hooks.example.com/in", EventFilter: []string{"deploy_started"}})
	if err != nil {
		t.Fatal(err)
	}
	if created.UUID != "w-1" || created.Secret != "whsec_example" || !created.Enabled {
		t.Errorf("unexpected webhook: %+v", created)
	}
}

func TestUpdateWebhookSendsEmptyAppUUIDs(t *testing.T) {
	var body string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.Write([]byte(`{"uuid":"w-1"}`))
	})
	empty := []string{}
	off := false
	if _, err := client.Webhooks.Update(context.Background(), "w-1", &UpdateWebhookRequest{AppUUIDs: &empty, Enabled: &off}); err != nil {
		t.Fatal(err)
	}
	if body != `{"app_uuids":[],"enabled":false}` {
		t.Errorf("body = %s", body)
	}
}
