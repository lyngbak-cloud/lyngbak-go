package miget

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestOperationCoverage verifies that every operationId in the OpenAPI spec
// has a corresponding implementation. Each SDK method carries a marker
// comment of the form:
//
//	// Operation: <operationId> (<METHOD> <path>)
//
// which this test cross-checks against openapi/openapi.json.
func TestOperationCoverage(t *testing.T) {
	specData, err := os.ReadFile(filepath.Join("openapi", "openapi.json"))
	if err != nil {
		t.Fatalf("reading spec: %v", err)
	}
	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string `json:"operationId"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(specData, &spec); err != nil {
		t.Fatalf("parsing spec: %v", err)
	}

	methods := map[string]bool{"get": true, "post": true, "put": true, "patch": true, "delete": true, "head": true, "options": true}
	var opIDs []string
	for _, ops := range spec.Paths {
		for method, op := range ops {
			if methods[strings.ToLower(method)] && op.OperationID != "" {
				opIDs = append(opIDs, op.OperationID)
			}
		}
	}
	if len(opIDs) != 115 {
		t.Errorf("spec contains %d operations, expected 115", len(opIDs))
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var source strings.Builder
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		source.Write(data)
		source.WriteByte('\n')
	}
	src := source.String()

	var missing []string
	for _, id := range opIDs {
		if !strings.Contains(src, "Operation: "+id+" (") {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		t.Errorf("unimplemented operations (%d): %s", len(missing), strings.Join(missing, ", "))
	}
}
