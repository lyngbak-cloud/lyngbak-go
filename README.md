# lyngbak-go

The official Go SDK for the [Miget](https://miget.com) PaaS REST API.

> Status: Generated from the published OpenAPI spec (fetched 2026-06-13); not yet validated against the live API.

- Zero third-party dependencies (standard library only)
- Covers all 115 documented API operations
- Context-aware: every method takes a `context.Context`
- Typed errors with `errors.As` support
- Server-Sent Events streaming for live deployment logs

## Install

```sh
go get github.com/lyngbak-cloud/lyngbak-go
```

Requires Go 1.22 or newer.

## Authentication

The client authenticates with a bearer token, sent as `Authorization: Bearer <token>`. Two kinds of tokens are accepted:

- API tokens (`miget_live_...`) created in the Miget dashboard
- Short-lived JWTs obtained from `client.Auth.SignIn`

```go
import "github.com/lyngbak-cloud/lyngbak-go"

// Explicit token:
client := lyngbak.New("miget_live_...")

// Or from the MIGET_API_TOKEN environment variable:
client, err := lyngbak.NewFromEnv()
```

Options:

```go
client := lyngbak.New(token,
    lyngbak.WithBaseURL("https://app.miget.com"), // default
    lyngbak.WithWorkspaceID("ws-uuid"),           // sets X-Workspace-Id on every request
    lyngbak.WithHTTPClient(httpClient),           // custom *http.Client
    lyngbak.WithUserAgent("my-tool/1.0"),         // custom User-Agent
)
```

Email/password sign-in (returns JWTs):

```go
tokens, err := client.Auth.SignIn(ctx, "user@example.com", "password")
// tokens.AccessToken (30 min), tokens.RefreshToken (7 days)
fresh, err := client.Auth.RefreshToken(ctx, tokens.RefreshToken)
```

## Workspaces

If your account belongs to multiple workspaces, select one with
`lyngbak.WithWorkspaceID(...)`. Without it the API uses your default workspace.

## Quickstart

Provision a resource, create a project and an app, deploy it, and poll the
deployment:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/lyngbak-cloud/lyngbak-go"
)

func main() {
    ctx := context.Background()
    client, err := lyngbak.NewFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 1. Provision a compute resource (a "Miget").
    res, err := client.Resources.Create(ctx, &lyngbak.CreateResourceRequest{
        PlanType:     "dev",
        PlanCodeName: "starter",
        RegionCode:   "eu-central-1",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 2. Create a project.
    project, err := client.Projects.Create(ctx, &lyngbak.CreateProjectRequest{
        Name: "demo",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 3. Create an app on the resource.
    app, err := client.Apps.Create(ctx, &lyngbak.CreateAppRequest{
        Label:            "Demo Web",
        Name:             "demo-web",
        ProjectID:        project.UUID,
        ResourceID:       res.UUID,
        Builder:          "auto",
        DeploymentMethod: "container_registry",
        DeploymentConfig: &lyngbak.AppDeploymentConfig{
            ImageURL: "docker.io/library/nginx",
            Tag:      "latest",
        },
        AppVars: []lyngbak.AppVarAttribute{
            {Key: "RAILS_ENV", Value: "production"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // 4. Trigger a deployment.
    if _, err := client.Apps.Deploy(ctx, app.UUID, nil); err != nil {
        log.Fatal(err)
    }

    // 5. Poll the latest deployment until it finishes.
    for {
        deployments, err := client.Deployments.List(ctx, app.UUID, nil)
        if err != nil {
            log.Fatal(err)
        }
        if len(deployments) > 0 {
            d := deployments[0]
            fmt.Println("deployment", d.UUID, "state:", d.State)
            if d.State == "completed" || d.State == "failed" || d.State == "cancelled" {
                break
            }
        }
        time.Sleep(5 * time.Second)
    }
}
```

## Error handling

Every non-2xx response is returned as a `*lyngbak.APIError` carrying the HTTP
status code and the parsed error payload. Both error body shapes used by the
API (`{"error": "..."}` and `{"errors": [...]}`) are parsed; unparseable
bodies are preserved raw.

```go
app, err := client.Apps.Get(ctx, uuid)
if err != nil {
    if lyngbak.IsNotFound(err) {
        // 404: app does not exist
    }
    var apiErr *lyngbak.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.StatusCode, apiErr.Message, apiErr.Errors)
    }
    return err
}
```

## Deployment logs

Stored logs (plain text):

```go
logs, err := client.Deployments.Logs(ctx, appUUID, deploymentID)
fmt.Print(logs)
```

Live logs over Server-Sent Events; the callback is invoked once per `data:`
line until the server closes the stream, the context is cancelled, or the
callback returns an error:

```go
ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
defer cancel()

err := client.Deployments.StreamLogs(ctx, appUUID, deploymentID, func(line string) error {
    fmt.Println(line)
    return nil // return an error to stop streaming early
})
```

Note: a custom `http.Client` with a `Timeout` set will also cut off
long-lived streams; prefer context deadlines.

## Optional fields

Optional boolean and numeric request fields are pointers so that `false`/`0`
can be sent explicitly. Use the `lyngbak.Ptr` helper:

```go
_, err := client.Apps.Update(ctx, uuid, &lyngbak.UpdateAppRequest{
    Maintenance: lyngbak.Ptr(true),
    RAMSize:     lyngbak.Ptr(512.0),
})
```

Optional string fields use `omitempty`; an empty string is not sent.

## Services and methods

| Service | Methods |
|---|---|
| `client.Apps` | `List`, `Create`, `Get`, `Update`, `Delete`, `UpdateState`, `UpdateSecurity`, `Clone`, `Deploy`, `UpdateDeployment`, `Activity`, `UpdateScalingProfile`, `UpdateHealthChecks`, `SetProjectSecretFilesEnabled` |
| `client.Projects` | `List`, `Create`, `Get`, `Update`, `Delete`, `ListApps`, `AssignResource`, `UnassignResource`, `CreateRestriction`, `DeleteRestriction` |
| `client.ProjectVars` | `List`, `Create`, `Update`, `Delete` |
| `client.Resources` | `List`, `Create`, `Get`, `Update`, `Delete`, `UpdateLabels` |
| `client.Services` | `List`, `Create`, `Get`, `Update`, `Delete`, `UpdateState`, `MountApp`, `UnmountApp`, `RotatePassword`, `Backups`, `RestoreBackup`, `ResetDatabase`, `CreateReplica`, `PromoteReplica`, `PromoteExternalReplica` |
| `client.Buckets` | `List`, `Create`, `Get`, `Update`, `Delete`, `UpdatePolicy`, `DeletePolicy`, `UpdateACL`, `DeleteACL` |
| `client.BucketObjects` | `List`, `GenerateUploadURL`, `GenerateDownloadURL`, `CreateFolder`, `Delete`, `Rename` |
| `client.Addons` | `List`, `Create`, `Get`, `Update`, `Delete`, `UpdateState`, `RotatePassword`, `Backups`, `RestoreBackup`, `ResetDatabase`, `CreateReplica`, `PromoteReplica`, `PromoteExternalReplica` |
| `client.Deployments` | `List`, `Get`, `Logs`, `StreamLogs`, `Cancel`, `Rollback` |
| `client.Domains` | `List`, `Create`, `Get`, `Update`, `Delete`, `Verify` |
| `client.Vars` | `List`, `Create`, `Update`, `Delete`, `SetProjectVariablesEnabled` |
| `client.Cronjobs` | `List`, `Create`, `Get`, `Update`, `Delete`, `StreamLogs` |
| `client.Ports` | `List`, `Create`, `Get`, `Delete`, `ExposePublicly`, `MakePrivate` |
| `client.Auth` | `SignIn`, `RefreshToken` |
| `client.Users` | `Me` |
| `client.SSHKeys` | `List`, `Create`, `Get`, `Delete` |
| `client.Credits` | `Get` |
| `client.Plans` | `List` |
| `client.Regions` | `List` |
| `client.Components` | `List` |
| `client.RegistryCredentials` | `List`, `Create`, `Get`, `Update`, `Delete` |

Each method's doc comment carries the OpenAPI `operationId` and HTTP route,
and `coverage_test.go` cross-checks the implementation against
`openapi/openapi.json` (115/115 operations).

## License

MIT, see [LICENSE](LICENSE).
