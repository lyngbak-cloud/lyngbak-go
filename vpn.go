package lyngbak

import (
	"context"
	"net/http"
)

type VPNService struct {
	client *Client
}

func (s *VPNService) ListGateways(ctx context.Context, vpcUUID string) ([]VPNGateway, error) {
	var out []VPNGateway
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "vpn_gateways"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) CreateGateway(ctx context.Context, vpcUUID string, req *CreateVPNGatewayRequest) (*VPNGateway, error) {
	out := new(VPNGateway)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "vpn_gateways"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) GetGateway(ctx context.Context, vpcUUID, childUUID string) (*VPNGateway, error) {
	out := new(VPNGateway)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "vpn_gateways", childUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) DeleteGateway(ctx context.Context, vpcUUID, childUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "vpn_gateways", childUUID), nil, nil, nil)
}

func (s *VPNService) ListUsers(ctx context.Context, vpcUUID string) ([]VPNUser, error) {
	var out []VPNUser
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "vpn_users"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) CreateUser(ctx context.Context, vpcUUID string, req *CreateVPNUserRequest) (*VPNUser, error) {
	out := new(VPNUser)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "vpn_users"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) GetUser(ctx context.Context, vpcUUID, childUUID string) (*VPNUser, error) {
	out := new(VPNUser)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "vpn_users", childUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) DeleteUser(ctx context.Context, vpcUUID, childUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "vpn_users", childUUID), nil, nil, nil)
}

func (s *VPNService) GetUserConfig(ctx context.Context, vpcUUID, userUUID string) (*VPNUserConfig, error) {
	out := new(VPNUserConfig)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "vpn_users", userUUID, "config"), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) ListPeerRoutes(ctx context.Context, vpcUUID string) ([]VPCPeerRoute, error) {
	var out []VPCPeerRoute
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "peer_routes"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) CreatePeerRoute(ctx context.Context, vpcUUID string, req *CreateVPCPeerRouteRequest) (*VPCPeerRoute, error) {
	out := new(VPCPeerRoute)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "peer_routes"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) GetPeerRoute(ctx context.Context, vpcUUID, childUUID string) (*VPCPeerRoute, error) {
	out := new(VPCPeerRoute)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "peer_routes", childUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) DeletePeerRoute(ctx context.Context, vpcUUID, childUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "peer_routes", childUUID), nil, nil, nil)
}

func (s *VPNService) ApplyPeerRoutes(ctx context.Context, vpcUUID string) ([]VPCPeerRoute, error) {
	var out []VPCPeerRoute
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "peer_routes", "apply"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) ListSiteConnections(ctx context.Context, vpcUUID string) ([]SiteConnection, error) {
	var out []SiteConnection
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "site_connections"), nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) CreateSiteConnection(ctx context.Context, vpcUUID string, req *CreateVPCSiteConnectionRequest) (*SiteConnection, error) {
	out := new(SiteConnection)
	if err := s.client.do(ctx, http.MethodPost, vpcPath(vpcUUID, "site_connections"), nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) GetSiteConnection(ctx context.Context, vpcUUID, childUUID string) (*SiteConnection, error) {
	out := new(SiteConnection)
	if err := s.client.do(ctx, http.MethodGet, vpcPath(vpcUUID, "site_connections", childUUID), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VPNService) DeleteSiteConnection(ctx context.Context, vpcUUID, childUUID string) error {
	return s.client.do(ctx, http.MethodDelete, vpcPath(vpcUUID, "site_connections", childUUID), nil, nil, nil)
}
