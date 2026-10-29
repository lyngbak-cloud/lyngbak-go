package lyngbak

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type StaticSitesService struct {
	client *Client
}

func sitePath(uuid string, rest ...string) string {
	p := "/api/v1/static"
	if uuid != "" {
		p += "/" + url.PathEscape(uuid)
	}
	for _, r := range rest {
		p += "/" + r
	}
	return p
}

func (s *StaticSitesService) List(ctx context.Context) ([]StaticSite, error) {
	var out []StaticSite
	if err := s.client.do(ctx, http.MethodGet, sitePath(""), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) Create(ctx context.Context, req *CreateStaticSiteRequest) (*StaticSite, error) {
	out := new(StaticSite)
	if err := s.client.do(ctx, http.MethodPost, sitePath(""), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) Get(ctx context.Context, uuid string) (*StaticSite, error) {
	out := new(StaticSite)
	if err := s.client.do(ctx, http.MethodGet, sitePath(uuid), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) Update(ctx context.Context, uuid string, req *UpdateStaticSiteRequest) (*StaticSite, error) {
	out := new(StaticSite)
	if err := s.client.do(ctx, http.MethodPut, sitePath(uuid), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) Delete(ctx context.Context, uuid string) error {
	return s.client.do(ctx, http.MethodDelete, sitePath(uuid), nil, nil, nil)
}

func (s *StaticSitesService) UpdateDeployment(ctx context.Context, uuid string, req *UpdateStaticSiteDeploymentRequest) (*StaticSite, error) {
	out := new(StaticSite)
	if err := s.client.do(ctx, http.MethodPut, sitePath(uuid, "deployment"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) Deploy(ctx context.Context, uuid string, req *DeployStaticSiteRequest) (*Message, error) {
	out := new(Message)
	path := sitePath(uuid, "deployments")
	if req == nil || req.Archive == nil {
		if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
			return nil, err
		}
		return out, nil
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	filename := req.Filename
	if filename == "" {
		filename = "site.zip"
	}
	part, err := mw.CreateFormFile("archive", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, req.Archive); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	r, err := s.client.newRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, err
	}
	body := buf.Bytes()
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	r.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	r.Header.Set("Content-Type", mw.FormDataContentType())
	if err := s.client.send(r, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) ListDeployments(ctx context.Context, uuid string, opts *DeploymentListOptions) ([]Deployment, error) {
	query := url.Values{}
	if opts != nil {
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.Period != "" {
			query.Set("period", opts.Period)
		}
	}
	var out []Deployment
	if err := s.client.do(ctx, http.MethodGet, sitePath(uuid, "deployments"), query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) CancelDeployment(ctx context.Context, uuid, deploymentID string) (*Message, error) {
	out := new(Message)
	path := sitePath(uuid, "deployments", url.PathEscape(deploymentID), "cancel")
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) RollbackDeployment(ctx context.Context, uuid, deploymentID string) (*Message, error) {
	out := new(Message)
	path := sitePath(uuid, "deployments", url.PathEscape(deploymentID), "rollback")
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) ListFiles(ctx context.Context, uuid string, opts *BucketObjectListOptions) (*BucketObjectList, error) {
	query := url.Values{}
	if opts != nil {
		if opts.Prefix != "" {
			query.Set("prefix", opts.Prefix)
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Cursor != "" {
			query.Set("cursor", opts.Cursor)
		}
	}
	out := new(BucketObjectList)
	if err := s.client.do(ctx, http.MethodGet, sitePath(uuid, "files"), query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) ListDomains(ctx context.Context, uuid string) ([]Domain, error) {
	var out []Domain
	if err := s.client.do(ctx, http.MethodGet, sitePath(uuid, "domains"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) CreateDomain(ctx context.Context, uuid, name string) (*Domain, error) {
	out := new(Domain)
	if err := s.client.do(ctx, http.MethodPost, sitePath(uuid, "domains"), nil, newDomainBody(name), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) GetDomain(ctx context.Context, uuid, domainUUID string) (*Domain, error) {
	out := new(Domain)
	if err := s.client.do(ctx, http.MethodGet, sitePath(uuid, "domains", url.PathEscape(domainUUID)), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) UpdateDomain(ctx context.Context, uuid, domainUUID, name string) (*Domain, error) {
	out := new(Domain)
	if err := s.client.do(ctx, http.MethodPut, sitePath(uuid, "domains", url.PathEscape(domainUUID)), nil, newDomainBody(name), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StaticSitesService) DeleteDomain(ctx context.Context, uuid, domainUUID string) error {
	return s.client.do(ctx, http.MethodDelete, sitePath(uuid, "domains", url.PathEscape(domainUUID)), nil, nil, nil)
}

func (s *StaticSitesService) VerifyDomain(ctx context.Context, uuid, domainUUID string) (*Domain, error) {
	out := new(Domain)
	if err := s.client.do(ctx, http.MethodPost, sitePath(uuid, "domains", url.PathEscape(domainUUID), "verify"), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
