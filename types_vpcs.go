package lyngbak

import "time"

type VPC struct {
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	Label           *string   `json:"label"`
	Status          string    `json:"status"`
	RegionID        int64     `json:"region_id"`
	RegionCode      string    `json:"region_code"`
	CIDRV4          string    `json:"cidr_v4"`
	CIDRV6          *string   `json:"cidr_v6"`
	ResolverV4      *string   `json:"resolver_v4"`
	VPNPoolV4       *string   `json:"vpn_pool_v4"`
	SubnetCount     int       `json:"subnet_count"`
	AttachmentCount int       `json:"attachment_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type VPCSubnet struct {
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Family          string    `json:"family"`
	CIDRV4          *string   `json:"cidr_v4"`
	CIDRV6          *string   `json:"cidr_v6"`
	GatewayV4       *string   `json:"gateway_v4"`
	GatewayV6       *string   `json:"gateway_v6"`
	Public          bool      `json:"public"`
	AttachmentCount int       `json:"attachment_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type VPCAttachment struct {
	UUID           string    `json:"uuid"`
	Status         string    `json:"status"`
	AttachableType string    `json:"attachable_type"`
	AttachableUUID *string   `json:"attachable_uuid"`
	WorkloadName   string    `json:"workload_name"`
	SubnetUUID     *string   `json:"subnet_uuid"`
	SubnetName     *string   `json:"subnet_name"`
	FQDN           *string   `json:"fqdn"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateVPCRequest struct {
	RegionID int64  `json:"region_id"`
	Name     string `json:"name,omitempty"`
	Label    string `json:"label,omitempty"`
	CIDRV4   string `json:"cidr_v4,omitempty"`
	CIDRV6   string `json:"cidr_v6,omitempty"`
}

type UpdateVPCRequest struct {
	Label string `json:"label"`
}

type CreateVPCSubnetRequest struct {
	Name   string `json:"name"`
	Family string `json:"family,omitempty"`
	CIDRV4 string `json:"cidr_v4,omitempty"`
	CIDRV6 string `json:"cidr_v6,omitempty"`
	Public *bool  `json:"public,omitempty"`
}

type CreateVPCAttachmentRequest struct {
	SubnetUUID     string `json:"subnet_uuid"`
	AttachableType string `json:"attachable_type"`
	AttachableUUID string `json:"attachable_uuid"`
}
