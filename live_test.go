//go:build live

package lyngbak

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	liveAPIURL = "https://dev.miget.com/api/v1"
	livePrefix = "lyngbak-go-live-"
)

func TestMain(m *testing.M) {
	if os.Getenv(EnvAPIURL) != liveAPIURL {
		fmt.Println("skipping live tests: MIGET_API_URL is not the dev API")
		os.Exit(0)
	}
	if os.Getenv(EnvAPIToken) == "" || os.Getenv(EnvWorkspaceID) == "" {
		fmt.Println("MIGET_API_TOKEN and MIGET_WORKSPACE_ID must be set")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func randSuffix() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func jsonKeys(raw json.RawMessage) map[string]json.RawMessage {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	return m
}

func assertModeled(t *testing.T, raw json.RawMessage, v any, optional ...string) {
	t.Helper()
	walkModeled(t, "", raw, reflect.TypeOf(v), optional)
}

func jsonShape(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	var sub map[string]json.RawMessage
	switch {
	case strings.HasPrefix(trimmed, "{"):
		_ = json.Unmarshal(raw, &sub)
	case strings.HasPrefix(trimmed, "["):
		var items []json.RawMessage
		_ = json.Unmarshal(raw, &items)
		if len(items) == 0 {
			return "array(empty)"
		}
		if strings.HasPrefix(strings.TrimSpace(string(items[0])), "{") {
			_ = json.Unmarshal(items[0], &sub)
			return "array of object" + keyList(sub)
		}
		return "array of scalar"
	case trimmed == "null":
		return "null"
	case strings.HasPrefix(trimmed, "\""):
		return "string"
	case trimmed == "true" || trimmed == "false":
		return "bool"
	default:
		return "number"
	}
	return "object" + keyList(sub)
}

func keyList(m map[string]json.RawMessage) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return "{" + strings.Join(keys, ",") + "}"
}

func walkModeled(t *testing.T, path string, raw json.RawMessage, typ reflect.Type, optional []string) {
	t.Helper()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return
	}
	switch typ.Kind() {
	case reflect.Slice:
		elem := typ.Elem()
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct || elem == reflect.TypeOf(time.Time{}) {
			return
		}
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return
		}
		for i, item := range items {
			walkModeled(t, fmt.Sprintf("%s[%d]", path, i), item, elem, optional)
		}
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			return
		}
		obj := jsonKeys(raw)
		if obj == nil {
			t.Errorf("%s: expected object, got %s", path, jsonShape(raw))
			return
		}
		declared := map[string]bool{}
		collectStruct(t, path, obj, typ, declared, optional)
		var extra []string
		for k := range obj {
			if !declared[k] {
				extra = append(extra, k+" ("+jsonShape(obj[k])+")")
			}
		}
		sort.Strings(extra)
		if len(extra) > 0 {
			t.Logf("unmodeled keys at %q: %s", path, strings.Join(extra, "; "))
		}
	}
}

func collectStruct(t *testing.T, path string, obj map[string]json.RawMessage, typ reflect.Type, declared map[string]bool, optional []string) {
	t.Helper()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if f.Anonymous && name == "" {
			ft := f.Type
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				collectStruct(t, path, obj, ft, declared, optional)
			}
			continue
		}
		if name == "" || name == "-" {
			continue
		}
		declared[name] = true
		val, ok := obj[name]
		if !ok {
			for _, o := range optional {
				if strings.HasPrefix(path+"."+name, o) {
					t.Logf("declared but omitted live: %s.%s", path, name)
					declared[name] = true
					goto next
				}
			}
			t.Errorf("%s.%s: declared by type but absent from live response", path, name)
		next:
			continue
		}
		walkModeled(t, path+"."+name, val, f.Type, optional)
	}
}

func liveRaw(t *testing.T, c *Client, method, path string, body any) json.RawMessage {
	t.Helper()
	var raw json.RawMessage
	if err := c.do(context.Background(), method, path, nil, body, &raw); err != nil {
		t.Fatalf("raw %s %s: %v", method, path, err)
	}
	return raw
}

func firstElement(t *testing.T, raw json.RawMessage) json.RawMessage {
	t.Helper()
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		t.Fatalf("expected array: %v", err)
	}
	if len(items) == 0 {
		return nil
	}
	return items[0]
}

func cleanupDelete(t *testing.T, what string, fn func() error) {
	t.Cleanup(func() {
		var err error
		for attempt := 0; attempt < 12; attempt++ {
			err = fn()
			var apiErr *APIError
			if err == nil {
				return
			}
			if errors.As(err, &apiErr) {
				code := apiErr.StatusCode
				if code == http.StatusNotFound {
					return
				}
				if code < 500 && code != http.StatusConflict && code != http.StatusUnprocessableEntity {
					break
				}
			}
			time.Sleep(10 * time.Second)
		}
		t.Errorf("cleanup %s: %v", what, err)
	})
}

func waitFor(t *testing.T, what string, timeout, every time.Duration, done func() (bool, error)) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		ok, err := done()
		if err != nil {
			t.Fatalf("waiting for %s: %v", what, err)
		}
		if ok {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(every)
	}
}

func retryCreate(t *testing.T, what string, timeout time.Duration, fn func() error) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		err := fn()
		if err == nil {
			return
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || time.Now().After(deadline) {
			t.Fatalf("%s: %v", what, err)
		}
		time.Sleep(10 * time.Second)
	}
}

func TestLive(t *testing.T) {
	top := t
	ctx := context.Background()
	c, err := NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	suffix := randSuffix()
	name := livePrefix + suffix

	var (
		regionID    int64
		projectID   string
		resourceID  string
		appUUID     string
		serviceUUID string
		vpcUUID     string
		subnetUUID  string
	)

	t.Run("catalog", func(t *testing.T) {
		regions, err := c.Regions.List(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range regions {
			if r.Code == "eu-east-1" {
				regionID = r.ID
			}
		}
		if len(regions) == 0 {
			t.Fatal("regions list is empty")
		}
		if regionID == 0 {
			t.Fatal("region eu-east-1 not found")
		}
		assertModeled(t, firstElement(t, liveRaw(t, c, http.MethodGet, "/api/v1/regions", nil)), Region{})

		plans, err := c.Plans.List(ctx)
		if err != nil || len(plans) == 0 {
			t.Fatalf("plans: %v (%d)", err, len(plans))
		}
		assertModeled(t, firstElement(t, liveRaw(t, c, http.MethodGet, "/api/v1/plans", nil)), Plan{})

		components, err := c.Components.List(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(components) == 0 {
			t.Fatal("components list is empty")
		}
		assertModeled(t, firstElement(t, liveRaw(t, c, http.MethodGet, "/api/v1/components", nil)), Component{})

		creds, err := c.GitCredentials.List(ctx)
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
			t.Logf("git credentials not readable with this token (403)")
		} else if err != nil {
			t.Fatal(err)
		} else if len(creds) == 0 {
			t.Skip("no git credentials in the dev workspace; GitCredential not checked")
		} else {
			assertModeled(t, firstElement(t, liveRaw(t, c, http.MethodGet, "/api/v1/git_credentials", nil)), GitCredential{})
		}
	})
	if regionID == 0 {
		t.FailNow()
	}

	t.Run("project_and_resource", func(t *testing.T) {
		p, err := c.Projects.Create(ctx, &CreateProjectRequest{Name: name, Description: "lyngbak-go live test"})
		if err != nil {
			t.Fatal(err)
		}
		projectID = p.UUID
		cleanupDelete(top, "project", func() error { return c.Projects.Delete(ctx, projectID) })
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/projects/"+projectID, nil), Project{})

		before, err := c.Resources.List(ctx)
		if err != nil {
			t.Fatal(err)
		}
		snapshot := time.Now()
		known := map[string]bool{}
		for _, r := range before {
			known[r.UUID] = true
		}
		r, err := c.Resources.Create(ctx, &CreateResourceRequest{PlanType: "dev", PlanCodeName: "hobby-1c-1g", RegionCode: "eu-east-1"})
		var pending *PaymentPendingError
		var checkout *CheckoutRequiredError
		switch {
		case errors.As(err, &checkout):
			t.Fatal("dev workspace needs a saved payment method for the paid plan")
		case errors.As(err, &pending):
			t.Logf("create returned payment pending; locating the new resource by listing")
			var found string
			deadline := time.Now().Add(2 * time.Minute)
			for {
				now, err := c.Resources.List(ctx)
				if err != nil {
					t.Fatal(err)
				}
				var matches []string
				for _, x := range now {
					if !known[x.UUID] && x.Plan.CodeName == "hobby-1c-1g" && x.RegionCode == "eu-east-1" && x.CreatedAt.After(snapshot) {
						matches = append(matches, x.UUID)
					}
				}
				if len(matches) > 1 {
					t.Fatalf("%d resources match the pending order; refusing to adopt any", len(matches))
				}
				if len(matches) == 1 {
					found = matches[0]
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("no resource matching the pending order appeared; nothing adopted")
				}
				time.Sleep(5 * time.Second)
			}
			resourceID = found
		case err != nil:
			t.Fatal(err)
		default:
			resourceID = r.UUID
		}
		cleanupDelete(top, "resource", func() error { return c.Resources.Delete(ctx, resourceID) })

		waitFor(t, "resource active", 5*time.Minute, 5*time.Second, func() (bool, error) {
			got, err := c.Resources.Get(ctx, resourceID)
			if err != nil {
				return false, err
			}
			return got.State == "active", nil
		})
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/resources/"+resourceID, nil), Resource{})
		if got, _ := c.Resources.Get(ctx, resourceID); got.RegionCode != "eu-east-1" {
			t.Errorf("region_code = %q", got.RegionCode)
		}
	})
	if resourceID == "" || projectID == "" {
		t.FailNow()
	}

	t.Run("app", func(t *testing.T) {
		app, err := c.Apps.Create(ctx, &CreateAppRequest{
			Label:            "Live " + suffix,
			Name:             "lgl" + suffix,
			ProjectID:        projectID,
			ResourceID:       resourceID,
			Builder:          "auto",
			RAMSize:          Ptr(128.0),
			CPUSize:          Ptr(0.1),
			DeploymentMethod: "git_push",
		})
		if err != nil {
			t.Fatal(err)
		}
		appUUID = app.UUID
		cleanupDelete(top, "app", func() error { return c.Apps.Delete(ctx, appUUID) })

		if _, err := c.Apps.UpdateHealthChecks(ctx, appUUID, &UpdateAppHealthChecksRequest{
			LivenessProbeEnabled:                     Ptr(true),
			LivenessProbePath:                        "/healthz",
			LivenessProbeInitialDelaySeconds:         Ptr(int64(15)),
			LivenessProbeTimeoutSeconds:              Ptr(int64(4)),
			LivenessProbePeriodSeconds:               Ptr(int64(20)),
			LivenessProbeFailureThreshold:            Ptr(int64(5)),
			ReadinessEmailFailureNotificationEnabled: Ptr(true),
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Apps.UpdateScalingProfile(ctx, appUUID, &UpdateAppScalingProfileRequest{Replicas: Ptr(int64(1)), AutoScalingEnabled: Ptr(false)}); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Apps.UpdateSecurity(ctx, appUUID, &UpdateAppSecurityRequest{
			BasicAuthEnabled:  Ptr(true),
			BasicAuthUsername: "live",
			BasicAuthPassword: "Live-" + suffix + "-pass",
		}); err != nil {
			t.Fatal(err)
		}

		got, err := c.Apps.Get(ctx, appUUID)
		if err != nil {
			t.Fatal(err)
		}
		raw := liveRaw(t, c, http.MethodGet, "/api/v1/apps/"+appUUID, nil)
		assertModeled(t, raw, App{}, ".deployment_config.")
		if got.HealthChecks == nil || got.HealthChecks.LivenessProbePath == nil || *got.HealthChecks.LivenessProbePath != "/healthz" {
			t.Errorf("health check path not read back (health_checks present: %v)", got.HealthChecks != nil)
		}
		if got.ScalingProfile == nil || got.ScalingProfile.Replicas != 1 || got.ScalingProfile.AutoScalingEnabled {
			t.Errorf("scaling profile not read back (present: %v)", got.ScalingProfile != nil)
		}
		if !got.BasicAuthEnabled || got.BasicAuthUsername == nil || *got.BasicAuthUsername != "live" {
			t.Errorf("basic auth not read back")
		}
		if strings.Contains(string(raw), "basic_auth_password") {
			t.Errorf("response carries basic_auth_password")
		}
		if got.RAMSize != 128 {
			t.Errorf("ram_size = %v", got.RAMSize)
		}
		if got.ResourceID != resourceID {
			t.Errorf("resource_id = %q", got.ResourceID)
		}
		obj := jsonKeys(raw)
		for _, k := range []string{"ram_size", "cpu_size"} {
			if string(obj[k]) == "null" {
				t.Errorf("%s is null", k)
			}
		}
		var miget map[string]json.RawMessage
		_ = json.Unmarshal(obj["miget"], &miget)
		for k, v := range miget {
			if strings.HasSuffix(k, "_size") && string(v) == "null" {
				t.Errorf("miget.%s is null", k)
			}
		}
		t.Logf("app addons: %s, cronjobs: %s", strings.TrimSpace(fmt.Sprint(len(got.Addons))), strings.TrimSpace(fmt.Sprint(len(got.Cronjobs))))
	})
	if appUUID == "" {
		t.FailNow()
	}

	t.Run("cronjob", func(t *testing.T) {
		cj, err := c.Cronjobs.Create(ctx, appUUID, &CreateAppCronjobRequest{
			Label:        "live",
			ScheduleType: "interval",
			IntervalType: "hourly",
			Command:      "echo live",
			Minute:       "0",
		})
		if err != nil {
			t.Fatal(err)
		}
		cleanupDelete(top, "cronjob", func() error { return c.Cronjobs.Delete(ctx, appUUID, cj.UUID) })
		raw := liveRaw(t, c, http.MethodGet, "/api/v1/apps/"+appUUID+"/cronjobs/"+cj.UUID, nil)
		assertModeled(t, raw, Cronjob{})
		appRaw := jsonKeys(liveRaw(t, c, http.MethodGet, "/api/v1/apps/"+appUUID, nil))
		var embedded []json.RawMessage
		if err := json.Unmarshal(appRaw["cronjobs"], &embedded); err != nil || len(embedded) == 0 {
			t.Fatalf("app.cronjobs: %v (%d)", err, len(embedded))
		}
		var keys []string
		for k := range jsonKeys(embedded[0]) {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		t.Logf("app.cronjobs[0] keys: %s", strings.Join(keys, ", "))
		assertModeled(t, appRaw["cronjobs"], []Addon{})
		app, err := c.Apps.Get(ctx, appUUID)
		if err != nil {
			t.Fatal(err)
		}
		if len(app.Cronjobs) != 1 {
			t.Errorf("app.Cronjobs = %d", len(app.Cronjobs))
		}
	})

	t.Run("secret_files_and_variables", func(t *testing.T) {
		sf, err := c.SecretFiles.Create(ctx, appUUID, &SecretFileRequest{Filename: "/app/config/live.json", Text: "one"})
		if err != nil {
			t.Fatal(err)
		}
		cleanupDelete(top, "app secret file", func() error { return c.SecretFiles.Delete(ctx, appUUID, sf.Name) })
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/apps/"+appUUID+"/secret_files/"+sf.Name, nil), SecretFile{})
		got, err := c.SecretFiles.Get(ctx, appUUID, sf.Name)
		if err != nil {
			t.Fatal(err)
		}
		if got.Filename != "/app/config/live.json" {
			t.Fatalf("secret file filename = %q", got.Filename)
		}
		if _, err := c.SecretFiles.Update(ctx, appUUID, sf.Name, &SecretFileRequest{Text: "two"}); err != nil {
			t.Fatal(err)
		}
		if err := c.SecretFiles.Delete(ctx, appUUID, sf.Name); err != nil {
			t.Fatal(err)
		}

		psf, err := c.ProjectSecretFiles.Create(ctx, projectID, &SecretFileRequest{Filename: "/app/config/project.json", Text: "one"})
		if err != nil {
			t.Fatal(err)
		}
		cleanupDelete(top, "project secret file", func() error { return c.ProjectSecretFiles.Delete(ctx, projectID, psf.Name) })
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/projects/"+projectID+"/secret_files/"+psf.Name, nil), SecretFile{})
		if _, err := c.ProjectSecretFiles.Get(ctx, projectID, psf.Name); err != nil {
			t.Fatal(err)
		}
		if err := c.ProjectSecretFiles.Delete(ctx, projectID, psf.Name); err != nil {
			t.Fatal(err)
		}

		if _, err := c.Vars.SetProjectVariablesEnabled(ctx, appUUID, true); err != nil {
			t.Fatal(err)
		}
		app, err := c.Apps.Get(ctx, appUUID)
		if err != nil {
			t.Fatal(err)
		}
		if !app.ProjectVariablesEnabled {
			t.Errorf("project_variables_enabled not read back")
		}
	})

	t.Run("service", func(t *testing.T) {
		svc, err := c.Services.Create(ctx, &CreateServiceRequest{
			ServiceType:     "postgres",
			ProjectID:       projectID,
			ResourceID:      resourceID,
			Label:           "Live db " + suffix,
			RAMSize:         Ptr(64.0),
			DiskSize:        Ptr(0.1),
			CPUSize:         Ptr(0.1),
			PostgresVersion: "16",
			PublicAccess:    "0",
			Instances:       Ptr(int64(1)),
		})
		if err != nil {
			t.Fatal(err)
		}
		serviceUUID = svc.UUID
		cleanupDelete(top, "service", func() error { return c.Services.Delete(ctx, serviceUUID) })
		raw := liveRaw(t, c, http.MethodGet, "/api/v1/services/"+serviceUUID, nil)
		assertModeled(t, raw, Service{})
		obj := jsonKeys(raw)
		for _, k := range []string{"ram_size", "disk_size", "cpu_size"} {
			if string(obj[k]) == "null" {
				t.Errorf("%s is null", k)
			}
		}
		got, err := c.Services.Get(ctx, serviceUUID)
		if err != nil {
			t.Fatal(err)
		}
		if got.RAMSize != 64 || got.DiskSize != 0.1 || got.CPUSize != 0.1 {
			t.Errorf("sizes: %v %v %v", got.RAMSize, got.DiskSize, got.CPUSize)
		}
		if got.PostgresVersion == nil || *got.PostgresVersion != "16" {
			t.Errorf("postgres_version = %v", got.PostgresVersion)
		}
		if got.ResourceID != resourceID {
			t.Errorf("resource_id = %q", got.ResourceID)
		}
	})

	t.Run("vpc", func(t *testing.T) {
		vpc, err := c.VPCs.Create(ctx, &CreateVPCRequest{RegionID: int(regionID), Name: name, Label: name})
		if err != nil {
			t.Fatal(err)
		}
		vpcUUID = vpc.UUID
		cleanupDelete(top, "vpc", func() error {
			subs, err := c.VPCs.ListSubnets(ctx, vpcUUID)
			if err != nil {
				return err
			}
			for _, s := range subs {
				_ = c.VPCs.DeleteSubnet(ctx, vpcUUID, s.UUID)
			}
			return c.VPCs.Delete(ctx, vpcUUID)
		})
		waitFor(t, "vpc active", 5*time.Minute, 5*time.Second, func() (bool, error) {
			got, err := c.VPCs.Get(ctx, vpcUUID)
			if err != nil {
				return false, err
			}
			return got.Status == "active", nil
		})
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/vpcs/"+vpcUUID, nil), VPC{})
		got, err := c.VPCs.Get(ctx, vpcUUID)
		if err != nil {
			t.Fatal(err)
		}
		prefix, err := netip.ParsePrefix(got.CIDRV4)
		if err != nil {
			t.Fatal(err)
		}
		a := prefix.Addr().As4()
		a[2], a[3] = 2, 0
		cidr := netip.PrefixFrom(netip.AddrFrom4(a), 24).String()

		var sub *VPCSubnet
		retryCreate(t, "create subnet", 2*time.Minute, func() error {
			var err error
			sub, err = c.VPCs.CreateSubnet(ctx, vpcUUID, &CreateVPCSubnetRequest{Name: "live" + suffix, Family: "ipv4", CIDRV4: cidr})
			return err
		})
		subnetUUID = sub.UUID
		cleanupDelete(top, "subnet", func() error { return c.VPCs.DeleteSubnet(ctx, vpcUUID, subnetUUID) })
		waitFor(t, "subnet active", 5*time.Minute, 5*time.Second, func() (bool, error) {
			s, err := c.VPCs.GetSubnet(ctx, vpcUUID, subnetUUID)
			if err != nil {
				return false, err
			}
			return s.Status == "active", nil
		})
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/vpcs/"+vpcUUID+"/subnets/"+subnetUUID, nil), VPCSubnet{})

		var att *VPCAttachment
		retryCreate(t, "create attachment", 2*time.Minute, func() error {
			var err error
			att, err = c.VPCs.CreateAttachment(ctx, vpcUUID, &CreateVPCAttachmentRequest{SubnetUUID: subnetUUID, AttachableType: "App", AttachableUUID: appUUID})
			return err
		})
		cleanupDelete(top, "attachment", func() error { return c.VPCs.DeleteAttachment(ctx, vpcUUID, att.UUID) })
		gotAtt, err := c.VPCs.GetAttachment(ctx, vpcUUID, att.UUID)
		if err != nil {
			t.Fatal(err)
		}
		if gotAtt.SubnetUUID == nil || *gotAtt.SubnetUUID != subnetUUID {
			t.Errorf("attachment subnet_uuid = %v", gotAtt.SubnetUUID)
		}
		attRaw := liveRaw(t, c, http.MethodGet, "/api/v1/vpcs/"+vpcUUID+"/attachments/"+att.UUID, nil)
		assertModeled(t, attRaw, VPCAttachment{})
		for _, k := range []string{"status", "fqdn", "subnet_name", "attachable_uuid"} {
			t.Logf("attachment.%s: %s", k, jsonShape(jsonKeys(attRaw)[k]))
		}
		subRaw := jsonKeys(liveRaw(t, c, http.MethodGet, "/api/v1/vpcs/"+vpcUUID+"/subnets/"+subnetUUID, nil))
		for _, k := range []string{"status", "family", "cidr_v6", "gateway_v6"} {
			t.Logf("subnet.%s: %s", k, jsonShape(subRaw[k]))
		}
	})

	t.Run("webhook", func(t *testing.T) {
		var apiErr *APIError
		wh, err := c.Webhooks.Create(ctx, &CreateWebhookRequest{
			Name:        name,
			URL:         "https://example.com/lyngbak-go-live",
			EventFilter: []string{"deploy_ended"},
		})
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
			t.Logf("webhooks not permitted for this token (403)")
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		cleanupDelete(top, "webhook", func() error { return c.Webhooks.Delete(ctx, wh.UUID) })
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/webhooks/"+wh.UUID, nil), Webhook{})
		got, err := c.Webhooks.Get(ctx, wh.UUID)
		if err != nil || got.URL != "https://example.com/lyngbak-go-live" {
			t.Fatalf("get: %v", err)
		}
		list, err := c.Webhooks.List(ctx)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, w := range list {
			found = found || w.UUID == wh.UUID
		}
		if !found {
			t.Errorf("webhook missing from list")
		}
	})

	t.Run("preview_environment_config", func(t *testing.T) {
		var apiErr *APIError
		if _, err := c.PreviewEnvironments.GetConfig(ctx, appUUID); errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			t.Logf("preview environment config not found for an app without a repository (404)")
			return
		} else if err != nil {
			t.Fatal(err)
		}
		assertModeled(t, liveRaw(t, c, http.MethodGet, "/api/v1/apps/"+appUUID+"/preview_environments/config", nil), PreviewEnvironmentConfig{})
	})
}
