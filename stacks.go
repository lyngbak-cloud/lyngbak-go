package lyngbak

import (
	"context"
	"net/http"
	"net/url"
)

type StacksService struct {
	client *Client
}

func stackPath(uuid string, rest ...string) string {
	p := "/api/v1/stacks"
	if uuid != "" {
		p += "/" + url.PathEscape(uuid)
	}
	for _, r := range rest {
		p += "/" + r
	}
	return p
}

func (s *StacksService) List(ctx context.Context) ([]Stack, error) {
	var out []Stack
	if err := s.client.do(ctx, http.MethodGet, stackPath(""), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) Create(ctx context.Context, req *CreateStackRequest) (*Stack, error) {
	out := new(Stack)
	if err := s.client.do(ctx, http.MethodPost, stackPath(""), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) Analyze(ctx context.Context, req *AnalyzeStackRequest) (*StackAnalysis, error) {
	out := new(StackAnalysis)
	if err := s.client.do(ctx, http.MethodPost, stackPath("", "analyze"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) Get(ctx context.Context, uuid string) (*Stack, error) {
	out := new(Stack)
	if err := s.client.do(ctx, http.MethodGet, stackPath(uuid), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) Update(ctx context.Context, uuid string, req *UpdateStackRequest) (*Stack, error) {
	out := new(Stack)
	if err := s.client.do(ctx, http.MethodPut, stackPath(uuid), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) Delete(ctx context.Context, uuid string) error {
	return s.client.do(ctx, http.MethodDelete, stackPath(uuid), nil, nil, nil)
}

func (s *StacksService) Deploy(ctx context.Context, uuid string, req *DeployStackRequest) (*Message, error) {
	out := new(Message)
	if err := s.client.do(ctx, http.MethodPost, stackPath(uuid, "deploy"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) UpdateDeployment(ctx context.Context, uuid string, req *UpdateStackDeploymentRequest) (*Stack, error) {
	out := new(Stack)
	if err := s.client.do(ctx, http.MethodPut, stackPath(uuid, "deployment"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) ListDeployments(ctx context.Context, uuid string) ([]StackDeployment, error) {
	var out []StackDeployment
	if err := s.client.do(ctx, http.MethodGet, stackPath(uuid, "deployments"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *StacksService) GetDeployment(ctx context.Context, uuid, id string) (*StackDeployment, error) {
	out := new(StackDeployment)
	if err := s.client.do(ctx, http.MethodGet, stackPath(uuid, "deployments", url.PathEscape(id)), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
