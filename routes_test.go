package lyngbak

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type sdkRoute struct {
	method, path string
	params       []string
	call         func(ctx context.Context, c *Client) error
}

var sdkRoutes = []sdkRoute{
	{method: "PUT", path: "/api/v1/apps/{uuid}/secret_files/project_files_enabled", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.SetProjectSecretFilesEnabled(ctx, "p1", true)
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/vars/project_variables_enabled", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Vars.SetProjectVariablesEnabled(ctx, "p1", true)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/cronjobs/{id}/stream_logs", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Cronjobs.StreamLogs(ctx, "p1", "p2", func(LogEvent) error { return nil })
	}},
	{method: "POST", path: "/api/v1/projects/{project_id}/resources", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.AssignResource(ctx, "p1", "p2")
		return err
	}},
	{method: "DELETE", path: "/api/v1/projects/{project_id}/resources/{resource_id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Projects.UnassignResource(ctx, "p1", "p2")
	}},
	{method: "POST", path: "/api/v1/projects/{project_id}/restrictions", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.CreateRestriction(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/projects/{project_id}/restrictions/{id}", params: []string{"p1", "2"}, call: func(ctx context.Context, c *Client) error {
		return c.Projects.DeleteRestriction(ctx, "p1", 2)
	}},
	{method: "GET", path: "/api/v1/apps", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/apps", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Apps.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Update(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/activity", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Activity(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/addons", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.Create(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/addons/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Addons.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/addons/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/addons/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.Update(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/addons/{id}/backups", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.Backups(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/create_replica", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.CreateReplica(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/promote_external", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.PromoteExternalReplica(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/promote_replica", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.PromoteReplica(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/reset_database", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.ResetDatabase(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/restore_backup", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.RestoreBackup(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/addons/{id}/rotate_password", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.RotatePassword(ctx, "p1", "p2")
		return err
	}},
	{method: "PATCH", path: "/api/v1/apps/{uuid}/addons/{id}/state", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Addons.UpdateState(ctx, "p1", "p2", "p3")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/clone", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Clone(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/cronjobs", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Cronjobs.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/cronjobs", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Cronjobs.Create(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/cronjobs/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Cronjobs.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/cronjobs/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Cronjobs.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/cronjobs/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Cronjobs.Update(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/deploy", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.Deploy(ctx, "p1", nil)
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/deployment", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.UpdateDeployment(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/deployments", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Deployments.List(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/deployments/{id}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Deployments.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/deployments/{id}/cancel", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Deployments.Cancel(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/deployments/{id}/logs", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Deployments.Logs(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/deployments/{id}/rollback", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Deployments.Rollback(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/deployments/{id}/stream_logs", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Deployments.StreamLogs(ctx, "p1", "p2", func(string) error { return nil })
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/domains", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/domains", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Create(ctx, "p1", "p2")
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/domains/{domain_uuid}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.Domains.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/domains/{domain_uuid}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/domains/{domain_uuid}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Update(ctx, "p1", "p2", "p3")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/domains/{domain_uuid}/verify", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Verify(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/health_checks", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.UpdateHealthChecks(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/ports", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Ports.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/ports", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Ports.Create(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/ports/{port_id}", params: []string{"p1", "1"}, call: func(ctx context.Context, c *Client) error {
		return c.Ports.Delete(ctx, "p1", 1)
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/ports/{port_id}", params: []string{"p1", "1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Ports.Get(ctx, "p1", 1)
		return err
	}},
	{method: "PATCH", path: "/api/v1/apps/{uuid}/ports/{port_id}/expose_publicly", params: []string{"p1", "1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Ports.ExposePublicly(ctx, "p1", 1)
		return err
	}},
	{method: "PATCH", path: "/api/v1/apps/{uuid}/ports/{port_id}/make_private", params: []string{"p1", "1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Ports.MakePrivate(ctx, "p1", 1)
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/scaling_profile", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.UpdateScalingProfile(ctx, "p1", nil)
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/security", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.UpdateSecurity(ctx, "p1", nil)
		return err
	}},
	{method: "PATCH", path: "/api/v1/apps/{uuid}/state", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Apps.UpdateState(ctx, "p1", "p2")
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Vars.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Vars.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Vars.Create(ctx, "p1", "p2", "p3")
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Vars.Update(ctx, "p1", nil)
	}},
	{method: "POST", path: "/api/v1/auth/refresh_token", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Auth.RefreshToken(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/auth/sign_in", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Auth.SignIn(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/buckets", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/buckets", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/buckets/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Buckets.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/buckets/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/buckets/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.Update(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/buckets/{uuid}/acl", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Buckets.DeleteACL(ctx, "p1")
	}},
	{method: "PUT", path: "/api/v1/buckets/{uuid}/acl", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.UpdateACL(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/buckets/{uuid}/objects/create_folder", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.BucketObjects.CreateFolder(ctx, "p1", "p2")
	}},
	{method: "POST", path: "/api/v1/buckets/{uuid}/objects/download_url", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.BucketObjects.GenerateDownloadURL(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/buckets/{uuid}/objects/list", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.BucketObjects.List(ctx, "p1", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/buckets/{uuid}/objects/upload_url", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.BucketObjects.GenerateUploadURL(ctx, "p1", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/buckets/{uuid}/objects/{key}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.BucketObjects.Delete(ctx, "p1", "p2")
	}},
	{method: "PUT", path: "/api/v1/buckets/{uuid}/objects/{key}/rename", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.BucketObjects.Rename(ctx, "p1", "p2", "p3")
	}},
	{method: "DELETE", path: "/api/v1/buckets/{uuid}/policy", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Buckets.DeletePolicy(ctx, "p1")
	}},
	{method: "PUT", path: "/api/v1/buckets/{uuid}/policy", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Buckets.UpdatePolicy(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/components", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Components.List(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/container_registry_credentials", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.RegistryCredentials.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/container_registry_credentials", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.RegistryCredentials.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/container_registry_credentials/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.RegistryCredentials.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/container_registry_credentials/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.RegistryCredentials.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/container_registry_credentials/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.RegistryCredentials.Update(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/plans", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Plans.List(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/projects", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/projects", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/projects/{project_id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Projects.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/projects/{project_id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/projects/{project_id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.Update(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/projects/{project_id}/apps", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Projects.ListApps(ctx, "p1")
		return err
	}},
	{method: "DELETE", path: "/api/v1/projects/{project_id}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.ProjectVars.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/projects/{project_id}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.ProjectVars.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/projects/{project_id}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.ProjectVars.Create(ctx, "p1", "p2", "p3")
	}},
	{method: "PUT", path: "/api/v1/projects/{project_id}/vars", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.ProjectVars.Update(ctx, "p1", "p2", "p3")
	}},
	{method: "GET", path: "/api/v1/regions", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Regions.List(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/resources", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Resources.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/resources", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Resources.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/resources/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Resources.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/resources/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Resources.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/resources/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Resources.Update(ctx, "p1", nil)
		return err
	}},
	{method: "PATCH", path: "/api/v1/resources/{uuid}/labels", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Resources.UpdateLabels(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/services", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/services", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.Create(ctx, nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/services/{id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		return c.Services.Delete(ctx, "p1")
	}},
	{method: "GET", path: "/api/v1/services/{id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.Get(ctx, "p1")
		return err
	}},
	{method: "PUT", path: "/api/v1/services/{id}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.Update(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/services/{id}/backups", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.Backups(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/create_replica", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.CreateReplica(ctx, "p1", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/mount_app", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.MountApp(ctx, "p1", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/promote_external", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.PromoteExternalReplica(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/promote_replica", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.PromoteReplica(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/reset_database", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.ResetDatabase(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/restore_backup", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.RestoreBackup(ctx, "p1", nil)
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/rotate_password", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.RotatePassword(ctx, "p1")
		return err
	}},
	{method: "PATCH", path: "/api/v1/services/{id}/state", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.UpdateState(ctx, "p1", "p2")
		return err
	}},
	{method: "POST", path: "/api/v1/services/{id}/unmount_app", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Services.UnmountApp(ctx, "p1", "p2")
		return err
	}},
	{method: "GET", path: "/api/v1/users/me", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.Me(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/users/me/credits", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.Credits.Get(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/users/me/ssh_keys", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SSHKeys.List(ctx)
		return err
	}},
	{method: "POST", path: "/api/v1/users/me/ssh_keys", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SSHKeys.Create(ctx, "p1")
		return err
	}},
	{method: "DELETE", path: "/api/v1/users/me/ssh_keys/{id}", params: []string{"1"}, call: func(ctx context.Context, c *Client) error {
		return c.SSHKeys.Delete(ctx, 1)
	}},
	{method: "GET", path: "/api/v1/users/me/ssh_keys/{id}", params: []string{"1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SSHKeys.Get(ctx, 1)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/secret_files", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecretFiles.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/apps/{uuid}/secret_files", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecretFiles.Create(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/apps/{uuid}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecretFiles.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/apps/{uuid}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecretFiles.Update(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/apps/{uuid}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.SecretFiles.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/projects/{project_id}/secret_files", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.ProjectSecretFiles.List(ctx, "p1")
		return err
	}},
	{method: "POST", path: "/api/v1/projects/{project_id}/secret_files", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.ProjectSecretFiles.Create(ctx, "p1", nil)
		return err
	}},
	{method: "GET", path: "/api/v1/projects/{project_id}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.ProjectSecretFiles.Get(ctx, "p1", "p2")
		return err
	}},
	{method: "PUT", path: "/api/v1/projects/{project_id}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.ProjectSecretFiles.Update(ctx, "p1", "p2", nil)
		return err
	}},
	{method: "DELETE", path: "/api/v1/projects/{project_id}/secret_files/{name}", params: []string{"p1", "p2"}, call: func(ctx context.Context, c *Client) error {
		return c.ProjectSecretFiles.Delete(ctx, "p1", "p2")
	}},
	{method: "GET", path: "/api/v1/git_credentials", params: []string{}, call: func(ctx context.Context, c *Client) error {
		_, err := c.GitCredentials.List(ctx)
		return err
	}},
	{method: "GET", path: "/api/v1/git_credentials/{uuid}", params: []string{"p1"}, call: func(ctx context.Context, c *Client) error {
		_, err := c.GitCredentials.Get(ctx, "p1")
		return err
	}},
}

var pendingRoutes = map[string]bool{
	"DELETE /api/v1/apps/{uuid}/preview_environments/{id}":          true,
	"DELETE /api/v1/stacks/{uuid}":                                  true,
	"DELETE /api/v1/static/{uuid}":                                  true,
	"DELETE /api/v1/static/{uuid}/domains/{domain_uuid}":            true,
	"DELETE /api/v1/vpcs/{uuid}":                                    true,
	"DELETE /api/v1/vpcs/{uuid}/attachments/{attachment_uuid}":      true,
	"DELETE /api/v1/vpcs/{uuid}/peer_routes/{route_uuid}":           true,
	"DELETE /api/v1/vpcs/{uuid}/site_connections/{connection_uuid}": true,
	"DELETE /api/v1/vpcs/{uuid}/subnets/{subnet_uuid}":              true,
	"DELETE /api/v1/vpcs/{uuid}/vpn_gateways/{gateway_uuid}":        true,
	"DELETE /api/v1/vpcs/{uuid}/vpn_users/{user_uuid}":              true,
	"DELETE /api/v1/webhooks/{uuid}":                                true,
	"GET /api/v1/apps/{uuid}/preview_environments":                  true,
	"GET /api/v1/apps/{uuid}/preview_environments/config":           true,
	"GET /api/v1/apps/{uuid}/preview_environments/{id}":             true,
	"GET /api/v1/stacks":                                            true,
	"GET /api/v1/stacks/{uuid}":                                     true,
	"GET /api/v1/stacks/{uuid}/deployments":                         true,
	"GET /api/v1/stacks/{uuid}/deployments/{id}":                    true,
	"GET /api/v1/static":                                            true,
	"GET /api/v1/static/{uuid}":                                     true,
	"GET /api/v1/static/{uuid}/deployments":                         true,
	"GET /api/v1/static/{uuid}/domains":                             true,
	"GET /api/v1/static/{uuid}/domains/{domain_uuid}":               true,
	"GET /api/v1/static/{uuid}/files":                               true,
	"GET /api/v1/vpcs":                                              true,
	"GET /api/v1/vpcs/{uuid}":                                       true,
	"GET /api/v1/vpcs/{uuid}/attachments":                           true,
	"GET /api/v1/vpcs/{uuid}/attachments/{attachment_uuid}":         true,
	"GET /api/v1/vpcs/{uuid}/peer_routes":                           true,
	"GET /api/v1/vpcs/{uuid}/peer_routes/{route_uuid}":              true,
	"GET /api/v1/vpcs/{uuid}/site_connections":                      true,
	"GET /api/v1/vpcs/{uuid}/site_connections/{connection_uuid}":    true,
	"GET /api/v1/vpcs/{uuid}/subnets":                               true,
	"GET /api/v1/vpcs/{uuid}/subnets/{subnet_uuid}":                 true,
	"GET /api/v1/vpcs/{uuid}/vpn_gateways":                          true,
	"GET /api/v1/vpcs/{uuid}/vpn_gateways/{gateway_uuid}":           true,
	"GET /api/v1/vpcs/{uuid}/vpn_users":                             true,
	"GET /api/v1/vpcs/{uuid}/vpn_users/{user_uuid}":                 true,
	"GET /api/v1/vpcs/{uuid}/vpn_users/{user_uuid}/config":          true,
	"GET /api/v1/webhooks":                                          true,
	"GET /api/v1/webhooks/{uuid}":                                   true,
	"GET /api/v1/webhooks/{uuid}/deliveries":                        true,
	"POST /api/v1/apps/{uuid}/preview_environments/{id}/redeploy":   true,
	"POST /api/v1/stacks":                                           true,
	"POST /api/v1/stacks/analyze":                                   true,
	"POST /api/v1/stacks/{uuid}/deploy":                             true,
	"POST /api/v1/static":                                           true,
	"POST /api/v1/static/{uuid}/deployments":                        true,
	"POST /api/v1/static/{uuid}/deployments/{id}/cancel":            true,
	"POST /api/v1/static/{uuid}/deployments/{id}/rollback":          true,
	"POST /api/v1/static/{uuid}/domains":                            true,
	"POST /api/v1/static/{uuid}/domains/{domain_uuid}/verify":       true,
	"POST /api/v1/vpcs":                                             true,
	"POST /api/v1/vpcs/{uuid}/attachments":                          true,
	"POST /api/v1/vpcs/{uuid}/peer_routes":                          true,
	"POST /api/v1/vpcs/{uuid}/peer_routes/apply":                    true,
	"POST /api/v1/vpcs/{uuid}/site_connections":                     true,
	"POST /api/v1/vpcs/{uuid}/subnets":                              true,
	"POST /api/v1/vpcs/{uuid}/vpn_gateways":                         true,
	"POST /api/v1/vpcs/{uuid}/vpn_users":                            true,
	"POST /api/v1/webhooks":                                         true,
	"POST /api/v1/webhooks/{uuid}/deliveries/{delivery_uuid}/retry": true,
	"POST /api/v1/webhooks/{uuid}/test":                             true,
	"PUT /api/v1/apps/{uuid}/preview_environments/config":           true,
	"PUT /api/v1/stacks/{uuid}":                                     true,
	"PUT /api/v1/stacks/{uuid}/deployment":                          true,
	"PUT /api/v1/static/{uuid}":                                     true,
	"PUT /api/v1/static/{uuid}/deployment":                          true,
	"PUT /api/v1/static/{uuid}/domains/{domain_uuid}":               true,
	"PUT /api/v1/vpcs/{uuid}":                                       true,
	"PUT /api/v1/webhooks/{uuid}":                                   true,
}

func specRouteKeys(t *testing.T) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("openapi", "openapi.json"))
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatal(err)
	}
	verbs := map[string]bool{"get": true, "post": true, "put": true, "patch": true, "delete": true}
	keys := map[string]bool{}
	for path, ops := range spec.Paths {
		for verb := range ops {
			if verbs[verb] {
				keys[strings.ToUpper(verb)+" "+path] = true
			}
		}
	}
	return keys
}

func TestRouteTableMatchesSpec(t *testing.T) {
	spec := specRouteKeys(t)
	if len(spec) != 206 {
		t.Errorf("spec has %d operations, want 206", len(spec))
	}
	seen := map[string]bool{}
	for _, r := range sdkRoutes {
		key := r.method + " " + r.path
		if seen[key] {
			t.Errorf("duplicate sdk route %s", key)
		}
		seen[key] = true
		if !spec[key] {
			t.Errorf("sdk route %s is not in the spec", key)
		}
		if pendingRoutes[key] {
			t.Errorf("%s is both implemented and pending", key)
		}
	}
	for key := range pendingRoutes {
		if !spec[key] {
			t.Errorf("pending route %s is not in the spec", key)
		}
	}
	for key := range spec {
		if !seen[key] && !pendingRoutes[key] {
			t.Errorf("spec route %s is neither implemented nor pending", key)
		}
	}
}

func TestRoutesHitDeclaredPath(t *testing.T) {
	for _, row := range sdkRoutes {
		t.Run(row.method+" "+row.path, func(t *testing.T) {
			var gotMethod, gotPath string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod, gotPath = r.Method, r.URL.EscapedPath()
				w.WriteHeader(http.StatusNoContent)
			}))
			defer srv.Close()
			c := New("t", WithBaseURL(srv.URL))
			_ = row.call(context.Background(), c)
			want := row.path
			for _, p := range row.params {
				i := strings.Index(want, "{")
				j := strings.Index(want, "}")
				want = want[:i] + p + want[j+1:]
			}
			if gotMethod != row.method {
				t.Errorf("method = %s, want %s", gotMethod, row.method)
			}
			if gotPath != want {
				t.Errorf("path = %s, want %s", gotPath, want)
			}
		})
	}
}
