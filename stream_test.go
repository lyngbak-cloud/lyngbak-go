package lyngbak

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDeploymentLogs(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/apps/app-1/deployments/d-1/logs" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Step 1/4 : FROM ruby\nStep 2/4 : COPY . .\n")
	})
	logs, err := client.Deployments.Logs(context.Background(), "app-1", "d-1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logs, "Step 2/4") {
		t.Errorf("unexpected logs: %q", logs)
	}
}

func TestDeploymentLogsError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"error":"logs not stored yet"}`)
	})
	_, err := client.Deployments.Logs(context.Background(), "app-1", "d-1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 {
		t.Fatalf("expected 503 *APIError, got %v", err)
	}
}

func TestStreamLogs(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/apps/app-1/deployments/d-1/stream_logs" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "text/event-stream" {
			t.Errorf("Accept = %q, want text/event-stream", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		// Mixed SSE input: comments, event fields, data with and without
		// the optional space, blank separators, and a final line without a
		// trailing newline.
		fmt.Fprint(w, ": keepalive comment\n")
		fmt.Fprint(w, "event: log\n")
		fmt.Fprint(w, "data: building image\n\n")
		flusher.Flush()
		fmt.Fprint(w, "data:no space here\n\n")
		fmt.Fprint(w, "id: 7\n")
		fmt.Fprint(w, "data:  double space\n\n")
		fmt.Fprint(w, "data: final line no newline")
		flusher.Flush()
	})

	var lines []string
	err := client.Deployments.StreamLogs(context.Background(), "app-1", "d-1", func(line string) error {
		lines = append(lines, line)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"building image", "no space here", " double space", "final line no newline"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %q, want %q", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestStreamLogsCallbackErrorStops(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: one\n\ndata: two\n\n")
	})
	stop := errors.New("stop")
	var count int
	err := client.Deployments.StreamLogs(context.Background(), "app-1", "d-1", func(line string) error {
		count++
		return stop
	})
	if !errors.Is(err, stop) {
		t.Fatalf("err = %v, want stop sentinel", err)
	}
	if count != 1 {
		t.Errorf("callback called %d times, want 1", count)
	}
}

func TestStreamLogsContextCancel(t *testing.T) {
	release := make(chan struct{})
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		fmt.Fprint(w, "data: first\n\n")
		w.(http.Flusher).Flush()
		<-release // hold the stream open until the test finishes
	})
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	err := client.Deployments.StreamLogs(ctx, "app-1", "d-1", func(line string) error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestStreamLogsHTTPError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"deployment not found"}`)
	})
	err := client.Deployments.StreamLogs(context.Background(), "app-1", "missing", func(line string) error { return nil })
	if !IsNotFound(err) {
		t.Fatalf("err = %v, want 404 *APIError", err)
	}
}

func TestSSEData(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"data: hello\n", "hello", true},
		{"data:hello\r\n", "hello", true},
		{"data:\n", "", true},
		{"data:  padded\n", " padded", true},
		{": comment\n", "", false},
		{"event: log\n", "", false},
		{"\n", "", false},
		{"id: 3\n", "", false},
	}
	for _, c := range cases {
		got, ok := sseData(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("sseData(%q) = (%q, %v), want (%q, %v)", c.in, got, ok, c.want, c.ok)
		}
	}
}
