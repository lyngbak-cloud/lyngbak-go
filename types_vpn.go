package lyngbak

import "time"

type VPNGateway struct {
	UUID                string    `json:"uuid"`
	Kind                string    `json:"kind"`
	Status              string    `json:"status"`
	Endpoint            *string   `json:"endpoint"`
	ExternalPort        *int      `json:"external_port"`
	PublicIP            *string   `json:"public_ip"`
	Organization        *string   `json:"organization"`
	ClientID            *string   `json:"client_id"`
	VPNUserCount        int       `json:"vpn_user_count"`
	SiteConnectionCount int       `json:"site_connection_count"`
	CreatedAt           time.Time `json:"created_at"`
}

type VPNUser struct {
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	PublicKey       *string   `json:"public_key"`
	AllocatedV4     *string   `json:"allocated_v4"`
	ConfigAvailable bool      `json:"config_available"`
	CreatedAt       time.Time `json:"created_at"`
}

type VPNUserConfig struct {
	Filename string `json:"filename"`
	Config   string `json:"config"`
}

type VPCPeerRoute struct {
	UUID        string    `json:"uuid"`
	CIDR        string    `json:"cidr"`
	GatewayKind string    `json:"gateway_kind"`
	Description *string   `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type SiteConnection struct {
	UUID           string    `json:"uuid"`
	Name           string    `json:"name"`
	PeerKind       string    `json:"peer_kind"`
	PeerASN        *int      `json:"peer_asn"`
	RemoteCIDRs    []string  `json:"remote_cidrs"`
	RemoteGateways []string  `json:"remote_gateways"`
	Status         string    `json:"status"`
	PublicIP       *string   `json:"public_ip"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateVPNGatewayRequest struct {
	Kind         string `json:"kind"`
	Organization string `json:"organization,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AuthKey      string `json:"auth_key,omitempty"`
}

type CreateVPNUserRequest struct {
	Name string `json:"name"`
}

type CreateVPCPeerRouteRequest struct {
	CIDR        string `json:"cidr"`
	GatewayKind string `json:"gateway_kind"`
	Description string `json:"description,omitempty"`
}

type SiteConnectionTunnel struct {
	RemoteGateway string `json:"remote_gateway"`
	PSK           string `json:"psk"`
}

type CreateVPCSiteConnectionRequest struct {
	Name        string                 `json:"name"`
	RemoteCIDRs []string               `json:"remote_cidrs"`
	Tunnels     []SiteConnectionTunnel `json:"tunnels"`
	PeerASN     *int                   `json:"peer_asn,omitempty"`
}
