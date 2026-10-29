package lyngbak

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestVPNUserConfigGone(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/vpn_users/u-1/config", 410, `{"error":"Gone"}`))
	_, err := client.VPN.GetUserConfig(context.Background(), "v-1", "u-1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 410 {
		t.Fatalf("expected 410 APIError, got %v", err)
	}
}

func TestVPNGatewayCloudflareFields(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/vpn_gateways/g-1", 200,
		`{"uuid":"g-1","kind":"cloudflare","status":"active","endpoint":null,"external_port":null,"public_ip":null,"organization":"example-team","client_id":"cid-example","vpn_user_count":0,"site_connection_count":0}`))
	gw, err := client.VPN.GetGateway(context.Background(), "v-1", "g-1")
	if err != nil {
		t.Fatal(err)
	}
	if gw.Organization == nil || *gw.Organization != "example-team" || gw.ClientID == nil || *gw.ClientID != "cid-example" {
		t.Errorf("unexpected gateway: %+v", gw)
	}
	if gw.Endpoint != nil || gw.ExternalPort != nil || gw.PublicIP != nil {
		t.Errorf("unexpected nullable fields: %+v", gw)
	}
}

func TestVPNGatewayWireguardNulls(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/vpn_gateways/g-2", 200,
		`{"uuid":"g-2","kind":"wireguard","status":"active","endpoint":"198.51.100.7","external_port":51820,"public_ip":null,"organization":null,"client_id":null,"vpn_user_count":1,"site_connection_count":0}`))
	gw, err := client.VPN.GetGateway(context.Background(), "v-1", "g-2")
	if err != nil {
		t.Fatal(err)
	}
	if gw.Organization != nil || gw.ClientID != nil || gw.ExternalPort == nil || *gw.ExternalPort != 51820 || gw.VPNUserCount != 1 {
		t.Errorf("unexpected gateway: %+v", gw)
	}
}

func TestApplyPeerRoutesReturnsSlice(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodPost, "/api/v1/vpcs/v-1/peer_routes/apply", 201,
		`[{"uuid":"r-1","cidr":"192.0.2.0/24","gateway_kind":"cloudflare","description":null,"status":"applied"},{"uuid":"r-2","cidr":"198.51.100.0/24","gateway_kind":"tailscale","description":"lab","status":"applied"}]`))
	routes, err := client.VPN.ApplyPeerRoutes(context.Background(), "v-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 2 || routes[0].CIDR != "192.0.2.0/24" || routes[0].Description != nil || routes[1].Description == nil || *routes[1].Description != "lab" {
		t.Errorf("unexpected routes: %+v", routes)
	}
}

func TestVPNUserPendingAndConfig(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/vpn_users/u-1", 200,
		`{"uuid":"u-1","name":"laptop","status":"pending","public_key":null,"allocated_v4":null,"config_available":false}`))
	u, err := client.VPN.GetUser(ctx, "v-1", "u-1")
	if err != nil {
		t.Fatal(err)
	}
	if u.PublicKey != nil || u.AllocatedV4 != nil || u.ConfigAvailable {
		t.Errorf("unexpected user: %+v", u)
	}
	client = newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/vpn_users/u-2/config", 200,
		`{"filename":"laptop.conf","config":"[Interface]\n"}`))
	cfg, err := client.VPN.GetUserConfig(ctx, "v-1", "u-2")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Filename != "laptop.conf" || cfg.Config != "[Interface]\n" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestSiteConnectionDecodes(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/site_connections/s-1", 200,
		`{"uuid":"s-1","name":"office","peer_kind":"aws","peer_asn":null,"remote_cidrs":["10.0.0.0/16"],"remote_gateways":["203.0.113.1"],"status":"active","public_ip":"192.0.2.9"}`))
	sc, err := client.VPN.GetSiteConnection(context.Background(), "v-1", "s-1")
	if err != nil {
		t.Fatal(err)
	}
	if sc.PeerASN != nil || len(sc.RemoteCIDRs) != 1 || len(sc.RemoteGateways) != 1 || sc.PublicIP == nil || *sc.PublicIP != "192.0.2.9" {
		t.Errorf("unexpected connection: %+v", sc)
	}
}

func TestVPNDeletesIgnoreBody(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		path string
		call func(c *Client) error
	}{
		{"/api/v1/vpcs/v-1/vpn_gateways/c-1", func(c *Client) error { return c.VPN.DeleteGateway(ctx, "v-1", "c-1") }},
		{"/api/v1/vpcs/v-1/vpn_users/c-1", func(c *Client) error { return c.VPN.DeleteUser(ctx, "v-1", "c-1") }},
		{"/api/v1/vpcs/v-1/peer_routes/c-1", func(c *Client) error { return c.VPN.DeletePeerRoute(ctx, "v-1", "c-1") }},
		{"/api/v1/vpcs/v-1/site_connections/c-1", func(c *Client) error { return c.VPN.DeleteSiteConnection(ctx, "v-1", "c-1") }},
	} {
		client := newTestClient(t, route(t, http.MethodDelete, tc.path, 200, `{"uuid":"c-1","status":"deleting"}`))
		if err := tc.call(client); err != nil {
			t.Errorf("%s: %v", tc.path, err)
		}
	}
}
