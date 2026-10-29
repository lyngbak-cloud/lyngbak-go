package lyngbak

import (
	"context"
	"net/http"
	"net/url"
)

type WebhooksService struct {
	client *Client
}

func webhookPath(uuid string, rest ...string) string {
	p := "/api/v1/webhooks"
	if uuid != "" {
		p += "/" + url.PathEscape(uuid)
	}
	for _, r := range rest {
		p += "/" + r
	}
	return p
}

func (s *WebhooksService) List(ctx context.Context) ([]Webhook, error) {
	var out []Webhook
	if err := s.client.do(ctx, http.MethodGet, webhookPath(""), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) Create(ctx context.Context, req *CreateWebhookRequest) (*WebhookCreated, error) {
	out := new(WebhookCreated)
	if err := s.client.do(ctx, http.MethodPost, webhookPath(""), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) Get(ctx context.Context, uuid string) (*Webhook, error) {
	out := new(Webhook)
	if err := s.client.do(ctx, http.MethodGet, webhookPath(uuid), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) Update(ctx context.Context, uuid string, req *UpdateWebhookRequest) (*Webhook, error) {
	out := new(Webhook)
	if err := s.client.do(ctx, http.MethodPut, webhookPath(uuid), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) Delete(ctx context.Context, uuid string) error {
	return s.client.do(ctx, http.MethodDelete, webhookPath(uuid), nil, nil, nil)
}

func (s *WebhooksService) Test(ctx context.Context, uuid string) (*WebhookDelivery, error) {
	out := new(WebhookDelivery)
	if err := s.client.do(ctx, http.MethodPost, webhookPath(uuid, "test"), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) ListDeliveries(ctx context.Context, uuid string) ([]WebhookDelivery, error) {
	var out []WebhookDelivery
	if err := s.client.do(ctx, http.MethodGet, webhookPath(uuid, "deliveries"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *WebhooksService) RetryDelivery(ctx context.Context, uuid, deliveryUUID string) (*WebhookDelivery, error) {
	out := new(WebhookDelivery)
	if err := s.client.do(ctx, http.MethodPost, webhookPath(uuid, "deliveries", url.PathEscape(deliveryUUID), "retry"), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
