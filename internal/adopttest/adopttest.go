// Package adopttest provides an in-memory Port API and plan builders for unit testing how resources handle
// objects that already exist in Port.
package adopttest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/require"
)

// collection is a kind of Port object served under /v1/<path>/{identifier}, wrapped in the response as <key>
type collection struct {
	path string
	key  string
}

var (
	blueprints   = collection{path: "blueprints", key: "blueprint"}
	pages        = collection{path: "pages", key: "page"}
	integrations = collection{path: "integration", key: "integration"}
)

// FakePort is an in-memory Port API serving the blueprint, page and integration endpoints, recording every
// write it receives
type FakePort struct {
	server *httptest.Server

	mu      sync.Mutex
	objects map[string]map[string]map[string]any
	writes  []string
	// appearOnCreate holds objects that appear when something tries to create them, as Port does when it
	// creates an object asynchronously after the resource checked that it did not exist
	appearOnCreate map[string]map[string]map[string]any
}

func NewFakePort(t testing.TB) *FakePort {
	t.Helper()
	f := &FakePort{objects: map[string]map[string]map[string]any{}, appearOnCreate: map[string]map[string]map[string]any{}}
	f.server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.server.Close)
	return f
}

// Client returns a Port client for the fake API with the given provider adoption policy
func (f *FakePort) Client(policy cli.AdoptionPolicy) *cli.PortClient {
	return &cli.PortClient{Client: resty.New().SetBaseURL(f.server.URL), AdoptExistingResources: policy}
}

func (f *FakePort) AddBlueprint(identifier string, blueprint map[string]any) {
	f.add(blueprints, identifier, blueprint)
}

func (f *FakePort) AddPage(identifier string, page map[string]any) {
	f.add(pages, identifier, page)
}

func (f *FakePort) AddIntegration(installationId string, integration map[string]any) {
	f.add(integrations, installationId, integration)
}

// PageAppearsOnCreate makes the page exist only once something tries to create it, so the create conflicts
func (f *FakePort) PageAppearsOnCreate(identifier string, page map[string]any) {
	f.appearOnCreateIn(pages, identifier, page)
}

// BlueprintAppearsOnCreate makes the blueprint exist only once something tries to create it, so the create conflicts
func (f *FakePort) BlueprintAppearsOnCreate(identifier string, blueprint map[string]any) {
	f.appearOnCreateIn(blueprints, identifier, blueprint)
}

func (f *FakePort) appearOnCreateIn(c collection, identifier string, object map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.appearOnCreate[c.path] == nil {
		f.appearOnCreate[c.path] = map[string]map[string]any{}
	}
	f.appearOnCreate[c.path][identifier] = object
}

func (f *FakePort) Blueprint(identifier string) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.objects[blueprints.path][identifier]
}

// Writes returns the non-GET requests received, as "METHOD /path"
func (f *FakePort) Writes() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.writes...)
}

func (f *FakePort) add(c collection, identifier string, object map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.objects[c.path] == nil {
		f.objects[c.path] = map[string]map[string]any{}
	}
	f.objects[c.path][identifier] = stamp(object)
}

func (f *FakePort) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/"), "/"), "/")
	var c collection
	switch parts[0] {
	case blueprints.path:
		c = blueprints
	case pages.path:
		c = pages
	case integrations.path:
		c = integrations
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"ok": false})
		return
	}
	if f.objects[c.path] == nil {
		f.objects[c.path] = map[string]map[string]any{}
	}
	objects := f.objects[c.path]

	if r.Method != http.MethodGet {
		f.writes = append(f.writes, r.Method+" "+r.URL.Path)
	}

	var body map[string]any
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}

	if len(parts) == 1 {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false})
			return
		}
		identifier, _ := body["identifier"].(string)
		if c == integrations {
			identifier, _ = body["installationId"].(string)
		}
		if pending, ok := f.appearOnCreate[c.path][identifier]; ok {
			objects[identifier] = stamp(pending)
			delete(f.appearOnCreate[c.path], identifier)
		}
		if _, exists := objects[identifier]; exists {
			writeJSON(w, http.StatusConflict, map[string]any{"ok": false, "error": "identifier_taken"})
			return
		}
		objects[identifier] = stamp(body)
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, c.key: objects[identifier]})
		return
	}

	identifier := parts[1]
	existing, exists := objects[identifier]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{"ok": false, "error": "not_found"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, c.key: existing})
	case http.MethodPut:
		objects[identifier] = stamp(body)
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, c.key: objects[identifier]})
	case http.MethodPatch:
		merged := map[string]any{}
		for k, v := range existing {
			merged[k] = v
		}
		for k, v := range body {
			if nested, ok := v.(map[string]any); ok {
				if existingNested, ok := merged[k].(map[string]any); ok {
					for nk, nv := range nested {
						existingNested[nk] = nv
					}
					continue
				}
			}
			merged[k] = v
		}
		objects[identifier] = stamp(merged)
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, c.key: merged})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false})
	}
}

// stamp sets the timestamps Port returns on every object
func stamp(object map[string]any) map[string]any {
	if object == nil {
		object = map[string]any{}
	}
	for _, field := range []string{"createdAt", "updatedAt"} {
		if _, ok := object[field]; !ok {
			object[field] = "2026-01-01T00:00:00.000Z"
		}
	}
	return object
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// ResourceSchema returns the schema of r
func ResourceSchema(t testing.TB, r resource.Resource) schema.Schema {
	t.Helper()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	return resp.Schema
}

// Plan builds a plan for s where every attribute is null except the given values
func Plan(t testing.TB, s schema.Schema, values map[string]tftypes.Value) tfsdk.Plan {
	t.Helper()
	objectType := s.Type().TerraformType(context.Background()).(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	for name, value := range values {
		_, known := objectType.AttributeTypes[name]
		require.True(t, known, "attribute %q is not in the schema", name)
		attributes[name] = value
	}
	return tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(objectType, attributes)}
}

// NullState returns an empty state for s, as for a resource being created
func NullState(s schema.Schema) tfsdk.State {
	return tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)}
}

func String(value string) tftypes.Value { return tftypes.NewValue(tftypes.String, value) }

func Bool(value bool) tftypes.Value { return tftypes.NewValue(tftypes.Bool, value) }

// Policy is a shorthand for a provider adoption policy
func Policy(kinds ...consts.ResourceKind) cli.AdoptionPolicy { return cli.NewAdoptionPolicy(kinds...) }

// ErrorSummaries returns the summaries of the error diagnostics
func ErrorSummaries(diags diag.Diagnostics) []string {
	var summaries []string
	for _, d := range diags.Errors() {
		summaries = append(summaries, d.Summary())
	}
	return summaries
}

// WarningSummaries returns the summaries of the warning diagnostics
func WarningSummaries(diags diag.Diagnostics) []string {
	var summaries []string
	for _, d := range diags.Warnings() {
		summaries = append(summaries, d.Summary())
	}
	return summaries
}

// State returns a state holding the same values as plan, as for a resource that already exists
func State(plan tfsdk.Plan) tfsdk.State {
	return tfsdk.State(plan)
}

// Config returns a configuration holding the same values as plan
func Config(plan tfsdk.Plan) tfsdk.Config {
	return tfsdk.Config(plan)
}
