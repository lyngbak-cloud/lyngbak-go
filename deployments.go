package lyngbak

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DeploymentsService handles the App Deployments endpoints (OpenAPI tag
// "App Deployments").
type DeploymentsService struct {
	client *Client
}

// List returns the deployments of an application. opts may be nil.
// Operation: listAppDeployments (GET /api/v1/apps/{uuid}/deployments)
func (s *DeploymentsService) List(ctx context.Context, appUUID string, opts *DeploymentListOptions) ([]Deployment, error) {
	query := url.Values{}
	if opts != nil {
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.Period != "" {
			query.Set("period", opts.Period)
		}
	}
	var out []Deployment
	path := fmt.Sprintf("/api/v1/apps/%s/deployments", url.PathEscape(appUUID))
	if err := s.client.do(ctx, http.MethodGet, path, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single deployment.
// Operation: getAppDeployment (GET /api/v1/apps/{uuid}/deployments/{id})
func (s *DeploymentsService) Get(ctx context.Context, appUUID, deploymentID string) (*Deployment, error) {
	out := new(Deployment)
	path := fmt.Sprintf("/api/v1/apps/%s/deployments/%s", url.PathEscape(appUUID), url.PathEscape(deploymentID))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Logs returns the stored build/deploy logs of a deployment as plain text.
// Operation: getAppDeploymentLogs (GET /api/v1/apps/{uuid}/deployments/{id}/logs)
func (s *DeploymentsService) Logs(ctx context.Context, appUUID, deploymentID string) (string, error) {
	path := fmt.Sprintf("/api/v1/apps/%s/deployments/%s/logs", url.PathEscape(appUUID), url.PathEscape(deploymentID))
	return s.client.doText(ctx, path)
}

// StreamLogs streams live deployment logs over Server-Sent Events
// (text/event-stream) and invokes fn once for every "data:" line received.
// It returns nil when the server closes the stream, ctx.Err() if the context
// is cancelled, or the error returned by fn (which also stops the stream).
// Operation: streamAppDeploymentLogs (GET /api/v1/apps/{uuid}/deployments/{id}/stream_logs)
func (s *DeploymentsService) StreamLogs(ctx context.Context, appUUID, deploymentID string, fn func(line string) error) error {
	path := fmt.Sprintf("/api/v1/apps/%s/deployments/%s/stream_logs", url.PathEscape(appUUID), url.PathEscape(deploymentID))
	return s.client.streamSSE(ctx, path, func(ev LogEvent) error { return fn(ev.Data) })
}

func (c *Client) streamSSE(ctx context.Context, path string, fn func(LogEvent) error) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return newAPIError(resp.StatusCode, body)
	}

	reader := bufio.NewReader(resp.Body)
	event := ""
	for {
		line, err := reader.ReadString('\n')
		// Process any data received before an error (including the final
		// unterminated line at EOF).
		if data, ok := sseData(line); ok {
			if fnErr := fn(LogEvent{Event: event, Data: data}); fnErr != nil {
				return fnErr
			}
		} else if name, ok := strings.CutPrefix(strings.TrimRight(line, "\r\n"), "event:"); ok {
			event = strings.TrimPrefix(name, " ")
		} else if strings.TrimRight(line, "\r\n") == "" {
			event = ""
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// sseData extracts the payload from an SSE "data:" line. It returns false
// for comments, other SSE fields (event:, id:, retry:) and blank lines.
func sseData(line string) (string, bool) {
	line = strings.TrimRight(line, "\r\n")
	if !strings.HasPrefix(line, "data:") {
		return "", false
	}
	data := strings.TrimPrefix(line, "data:")
	// Per the SSE spec, a single leading space after the colon is not part
	// of the payload.
	data = strings.TrimPrefix(data, " ")
	return data, true
}

// Cancel cancels a running deployment.
// Operation: cancelAppDeployment (POST /api/v1/apps/{uuid}/deployments/{id}/cancel)
func (s *DeploymentsService) Cancel(ctx context.Context, appUUID, deploymentID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/deployments/%s/cancel", url.PathEscape(appUUID), url.PathEscape(deploymentID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Rollback rolls the application back to this deployment.
// Operation: rollbackAppDeployment (POST /api/v1/apps/{uuid}/deployments/{id}/rollback)
func (s *DeploymentsService) Rollback(ctx context.Context, appUUID, deploymentID string) (*Message, error) {
	out := new(Message)
	path := fmt.Sprintf("/api/v1/apps/%s/deployments/%s/rollback", url.PathEscape(appUUID), url.PathEscape(deploymentID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
