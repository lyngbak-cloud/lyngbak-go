package lyngbak

import (
	"context"
	"net/http"
	"net/url"
)

type VPCsService struct {
	client *Client
}

func vpcPath(uuid string, rest ...string) string {
	p := "/api/v1/vpcs"
	if uuid != "" {
		p += "/" + url.PathEscape(uuid)
	}
	for _, r := range rest {
		p += "/" + url.PathEscape(r)
	}
	return p
}

func (s *VPCsService) List(ctx context.Context) ([]VPC, error) {
	var out []VPC
	if err := s.client.do(ctx, http.MethodGet, vpcPath(""), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) Create(ctx context.Context, req *CreateVPCRequest) (*VPC, error) {
	out := new(VPC)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(""), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) Get(ctx context.Context, uuid string) (*VPC, error) {
	out := new(VPC)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(uuid), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) Update(ctx context.Context, uuid string, req *UpdateVPCRequest) (*VPC, error) {
	out := new(VPC)
	if err := s.client.do(ctx, http.MethodPut, vpcPath(uuid), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) Delete(ctx context.Context, uuid string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(uuid), nil, nil, nil)
}

func (s *VPCsService) ListSubnets(ctx context.Context, vpcUUID string) ([]VPCSubnet, error) {
	var out []VPCSubnet
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "subnets"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) CreateSubnet(ctx context.Context, vpcUUID string, req *CreateVPCSubnetRequest) (*VPCSubnet, error) {
	out := new(VPCSubnet)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "subnets"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) GetSubnet(ctx context.Context, vpcUUID, subnetUUID string) (*VPCSubnet, error) {
	out := new(VPCSubnet)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "subnets", subnetUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) DeleteSubnet(ctx context.Context, vpcUUID, subnetUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "subnets", subnetUUID), nil, nil, nil)
}

func (s *VPCsService) ListAttachments(ctx context.Context, vpcUUID string) ([]VPCAttachment, error) {
	var out []VPCAttachment
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "attachments"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) CreateAttachment(ctx context.Context, vpcUUID string, req *CreateVPCAttachmentRequest) (*VPCAttachment, error) {
	out := new(VPCAttachment)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "attachments"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) GetAttachment(ctx context.Context, vpcUUID, attachmentUUID string) (*VPCAttachment, error) {
	out := new(VPCAttachment)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "attachments", attachmentUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPCsService) DeleteAttachment(ctx context.Context, vpcUUID, attachmentUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "attachments", attachmentUUID), nil, nil, nil)
}
