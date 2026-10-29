package lyngbak

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestVPCDecodesRegionID(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1", 200,
		`{"uuid":"v-1","region_id":1,"region_code":"eu-east-1","label":null,"cidr_v6":null,"resolver_v4":null,"vpn_pool_v4":null,"subnet_count":2,"attachment_count":3}`))
	vpc, err := client.VPCs.Get(context.Background(), "v-1")
	if err != nil {
		t.Fatal(err)
	}
	if vpc.RegionID != 1 || vpc.RegionCode != "eu-east-1" || vpc.SubnetCount != 2 || vpc.AttachmentCount != 3 {
		t.Errorf("unexpected vpc: %+v", vpc)
	}
	if vpc.Label != nil || vpc.CIDRV6 != nil || vpc.ResolverV4 != nil || vpc.VPNPoolV4 != nil {
		t.Errorf("unexpected nullable fields: %+v", vpc)
	}
}

func TestVPCDeleteIgnoresBody(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
	}{{204, ""}, {200, `{"uuid":"v-1","status":"deleting"}`}} {
		client := newTestClient(t, route(t, http.MethodDelete, "/api/v1/vpcs/v-1", tc.status, tc.body))
		if err := client.VPCs.Delete(context.Background(), "v-1"); err != nil {
			t.Errorf("status %d: %v", tc.status, err)
		}
	}
}

func TestVPCGetSubnetNotFound(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/subnets/s-1", 404, `{"error":"Not found"}`))
	_, err := client.VPCs.GetSubnet(context.Background(), "v-1", "s-1")
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false for %v", err)
	}
}

func TestVPCSubnetDecodes(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/subnets", 200,
		`[{"uuid":"s-1","name":"web","status":"active","family":"dual","cidr_v4":"10.0.1.0/24","cidr_v6":"2001:db8::/64","gateway_v4":"10.0.1.1","gateway_v6":null,"public":false,"attachment_count":1},{"uuid":"s-2","cidr_v4":null,"cidr_v6":null,"gateway_v4":null,"gateway_v6":null,"public":true}]`))
	list, err := client.VPCs.ListSubnets(context.Background(), "v-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].CIDRV4 == nil || *list[0].CIDRV4 != "10.0.1.0/24" || list[0].GatewayV6 != nil || list[0].Public || list[0].AttachmentCount != 1 {
		t.Errorf("unexpected subnets: %+v", list)
	}
	if list[1].CIDRV4 != nil || list[1].CIDRV6 != nil || list[1].GatewayV4 != nil || !list[1].Public {
		t.Errorf("unexpected nullable fields: %+v", list[1])
	}
}

func TestVPCAttachmentDecodesSubnetUUID(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/attachments/at-1", 200,
		`{"uuid":"at-1","status":"active","attachable_type":"App","attachable_uuid":"a-1","workload_name":"web","subnet_uuid":"s-1","subnet_name":"web","fqdn":"web.vpc.example.com"}`))
	at, err := client.VPCs.GetAttachment(context.Background(), "v-1", "at-1")
	if err != nil {
		t.Fatal(err)
	}
	if at.SubnetUUID == nil || *at.SubnetUUID != "s-1" || at.SubnetName == nil || at.AttachableUUID == nil || at.FQDN == nil || *at.FQDN != "web.vpc.example.com" || at.WorkloadName != "web" {
		t.Errorf("unexpected attachment: %+v", at)
	}
}

func TestVPCAttachmentNullableFields(t *testing.T) {
	client := newTestClient(t, route(t, http.MethodGet, "/api/v1/vpcs/v-1/attachments", 200,
		`[{"uuid":"at-2","attachable_uuid":null,"subnet_uuid":null,"subnet_name":null,"fqdn":null}]`))
	list, err := client.VPCs.ListAttachments(context.Background(), "v-1")
	if err != nil {
		t.Fatal(err)
	}
	if list[0].SubnetUUID != nil || list[0].SubnetName != nil || list[0].AttachableUUID != nil || list[0].FQDN != nil {
		t.Errorf("unexpected attachment: %+v", list[0])
	}
}

func TestVPCRequestBodies(t *testing.T) {
	var got map[string]any
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = nil
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{}`))
	})
	ctx := context.Background()
	if _, err := client.VPCs.Create(ctx, &CreateVPCRequest{RegionID: 1}); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got["region_id"] != float64(1) {
		t.Errorf("create body = %v", got)
	}
	if _, err := client.VPCs.Update(ctx, "v-1", &UpdateVPCRequest{Label: ""}); err != nil {
		t.Fatal(err)
	}
	if v, ok := got["label"]; !ok || v != "" {
		t.Errorf("update body = %v", got)
	}
	private := false
	if _, err := client.VPCs.CreateSubnet(ctx, "v-1", &CreateVPCSubnetRequest{Name: "web", Public: &private}); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got["public"] != false || got["name"] != "web" {
		t.Errorf("subnet body = %v", got)
	}
	if _, err := client.VPCs.CreateAttachment(ctx, "v-1", &CreateVPCAttachmentRequest{SubnetUUID: "s-1", AttachableType: "App", AttachableUUID: "a-1"}); err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got["subnet_uuid"] != "s-1" || got["attachable_type"] != "App" || got["attachable_uuid"] != "a-1" {
		t.Errorf("attachment body = %v", got)
	}
}
