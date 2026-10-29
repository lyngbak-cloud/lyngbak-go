package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DomainsService handles the App Domains endpoints (OpenAPI tag
// "App Domains").
type DomainsService struct {
	client *Client
}

// domainBody is the nested wire format used by the create/update domain
// endpoints: {"domain": {"name": "..."}}.
type domainBody struct {
	Domain struct {
		Name string `json:"name"`
	} `json:"domain"`
}

func newDomainBody(name string) domainBody {
	var b domainBody
	b.Domain.Name = name
	return b
}

// List returns all custom domains of an application.
// Operation: listAppDomains (GET /api/v1/apps/{uuid}/domains)
func (s *DomainsService) List(ctx context.Context, appUUID string) ([]Domain, error) {
	var out []Domain
	path := fmt.Sprintf("/api/v1/apps/%s/domains", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create adds a custom domain to an application. name is the fully
// qualified domain name, e.g. app.example.com.
// Operation: createAppDomain (POST /api/v1/apps/{uuid}/domains)
func (s *DomainsService) Create(ctx context.Context, appUUID, name string) (*Domain, error) {
	out := new(Domain)
	path := fmt.Sprintf("/api/v1/apps/%s/domains", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, newDomainBody(name), out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single custom domain.
// Operation: getAppDomain (GET /api/v1/apps/{uuid}/domains/{domain_uuid})
func (s *DomainsService) Get(ctx context.Context, appUUID, domainUUID string) (*Domain, error) {
	out := new(Domain)
	path := fmt.Sprintf("/api/v1/apps/%s/domains/%s", url.PathEscape(appUUID), url.PathEscape(domainUUID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a custom domain's name.
// Operation: updateAppDomain (PUT /api/v1/apps/{uuid}/domains/{domain_uuid})
func (s *DomainsService) Update(ctx context.Context, appUUID, domainUUID, name string) (*Domain, error) {
	out := new(Domain)
	path := fmt.Sprintf("/api/v1/apps/%s/domains/%s", url.PathEscape(appUUID), url.PathEscape(domainUUID))
	if err := s.client.do(ctx, http.MethodPut, path, nil, newDomainBody(name), out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes a custom domain from an application.
// Operation: deleteAppDomain (DELETE /api/v1/apps/{uuid}/domains/{domain_uuid})
func (s *DomainsService) Delete(ctx context.Context, appUUID, domainUUID string) error {
	path := fmt.Sprintf("/api/v1/apps/%s/domains/%s", url.PathEscape(appUUID), url.PathEscape(domainUUID))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Verify triggers DNS verification of a custom domain.
// Operation: verifyAppDomain (POST /api/v1/apps/{uuid}/domains/{domain_uuid}/verify)
func (s *DomainsService) Verify(ctx context.Context, appUUID, domainUUID string) (*Domain, error) {
	out := new(Domain)
	path := fmt.Sprintf("/api/v1/apps/%s/domains/%s/verify", url.PathEscape(appUUID), url.PathEscape(domainUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
