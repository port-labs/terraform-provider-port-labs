package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// workflowOnlyProvider serves just the workflow resource, so the protocol
// tests need no credentials and only reach the API through the client given.
type workflowOnlyProvider struct{ client *cli.PortClient }

func (p *workflowOnlyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "port"
}

func (p *workflowOnlyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *workflowOnlyProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.ResourceData = p.client
}

func (p *workflowOnlyProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewWorkflowResource}
}

func (p *workflowOnlyProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func newWorkflowRPCServer(t testing.TB, client *cli.PortClient) (tfprotov6.ProviderServer, tftypes.Type) {
	t.Helper()
	ctx := context.Background()
	server := providerserver.NewProtocol6(&workflowOnlyProvider{client: client})()

	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	requireNoErrorDiagnostics(t, schemas.Diagnostics)

	providerType := schemas.Provider.ValueType()
	configured, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		Config: dynamicValue(t, tftypes.NewValue(providerType, map[string]tftypes.Value{})),
	})
	require.NoError(t, err)
	requireNoErrorDiagnostics(t, configured.Diagnostics)

	return server, schemas.ResourceSchemas["port_workflow"].ValueType()
}

func requireNoErrorDiagnostics(t testing.TB, diagnostics []*tfprotov6.Diagnostic) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, "%s: %s", diagnostic.Summary, diagnostic.Detail)
	}
}

func rpcErrorSummaries(diagnostics []*tfprotov6.Diagnostic) []string {
	var summaries []string
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			summaries = append(summaries, diagnostic.Summary)
		}
	}
	return summaries
}

func dynamicValue(t testing.TB, value tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dynamic, err := tfprotov6.NewDynamicValue(value.Type(), value)
	require.NoError(t, err)
	return &dynamic
}

func valueAt(t testing.TB, value tftypes.Value, attributePath *tftypes.AttributePath) tftypes.Value {
	t.Helper()
	found, remaining, err := tftypes.WalkAttributePath(value, attributePath)
	require.NoError(t, err, "%s (unresolved %s)", attributePath, remaining)
	return found.(tftypes.Value)
}

// rpcWorkflowConfig builds a valid workflow of one self-serve trigger followed
// by webhook nodes, with a dataset-filtered entity input on the trigger.
func rpcWorkflowConfig(nodes int) configObject {
	service := configObject{
		"title":     "Service",
		"format":    "entity",
		"blueprint": "service",
		"dataset": configObject{
			"combinator": "and",
			"rules":      configList{configObject{"property": "$identifier", "operator": "=", "value_json": `"example"`}},
		},
	}
	nodeList := configList{configObject{
		"identifier": "trigger",
		"self_serve_trigger": configObject{
			"user_inputs": configObject{
				"user_properties": configObject{"string_props": configObject{"service": service}},
			},
		},
	}}
	for i := 1; i < nodes; i++ {
		nodeList = append(nodeList, configObject{
			"identifier": fmt.Sprintf("webhook_%d", i),
			"webhook":    configObject{"url": "https://example.invalid/hook"},
		})
	}
	return configObject{"identifier": "rpc_workflow", "title": "RPC regression", "node": nodeList}
}

func rpcNode(config configObject, index int) configObject {
	return config["node"].(configList)[index].(configObject)
}

func serviceUserProperty(config configObject) configObject {
	userInputs := rpcNode(config, 0)["self_serve_trigger"].(configObject)["user_inputs"].(configObject)
	return userInputs["user_properties"].(configObject)["string_props"].(configObject)["service"].(configObject)
}

func validateResourceConfig(t testing.TB, server tfprotov6.ProviderServer, config tftypes.Value) []*tfprotov6.Diagnostic {
	t.Helper()
	resp, err := server.ValidateResourceConfig(context.Background(), &tfprotov6.ValidateResourceConfigRequest{
		TypeName: "port_workflow",
		Config:   dynamicValue(t, config),
	})
	require.NoError(t, err)
	return resp.Diagnostics
}

func planResourceChange(t testing.TB, server tfprotov6.ProviderServer, config, prior tftypes.Value) *tfprotov6.PlanResourceChangeResponse {
	t.Helper()
	resp, err := server.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "port_workflow",
		Config:           dynamicValue(t, config),
		ProposedNewState: dynamicValue(t, config),
		PriorState:       dynamicValue(t, prior),
	})
	require.NoError(t, err)
	requireNoErrorDiagnostics(t, resp.Diagnostics)
	return resp
}

func TestWorkflowRPCDefaultsNullsAndUnknowns(t *testing.T) {
	server, resourceType := newWorkflowRPCServer(t, nil)
	config := configValue(t, resourceType, rpcWorkflowConfig(2))
	requireNoErrorDiagnostics(t, validateResourceConfig(t, server, config))

	planned, err := planResourceChange(t, server, config, tftypes.NewValue(resourceType, nil)).PlannedState.Unmarshal(resourceType)
	require.NoError(t, err)

	root := tftypes.NewAttributePath()
	trigger := root.WithAttributeName("node").WithElementKeyInt(0)
	webhook := root.WithAttributeName("node").WithElementKeyInt(1)
	defaults := map[*tftypes.AttributePath]any{
		root.WithAttributeName("allow_anyone_to_view_runs"):                            true,
		trigger.WithAttributeName("verbose"):                                           false,
		trigger.WithAttributeName("self_serve_trigger").WithAttributeName("published"): true,
		webhook.WithAttributeName("webhook").WithAttributeName("on_failure"):           "terminate",
	}
	for attributePath, want := range defaults {
		got := valueAt(t, planned, attributePath)
		assert.True(t, got.Equal(tftypes.NewValue(got.Type(), want)), "default %s = %s", attributePath, got)
	}

	numberProps := trigger.WithAttributeName("self_serve_trigger").WithAttributeName("user_inputs").
		WithAttributeName("user_properties").WithAttributeName("number_props")
	assert.True(t, valueAt(t, planned, numberProps).IsNull(), "absent numeric properties must remain null")
	assert.False(t, valueAt(t, planned, root.WithAttributeName("id")).IsKnown(), "new resource ID must remain unknown")

	for _, unknownAt := range []*tftypes.AttributePath{
		root.WithAttributeName("node"),
		webhook.WithAttributeName("webhook").WithAttributeName("url"),
		trigger.WithAttributeName("condition"),
	} {
		t.Run(unknownAt.String(), func(t *testing.T) {
			requireNoErrorDiagnostics(t, validateResourceConfig(t, server, withUnknownValue(t, config, unknownAt)))
		})
	}
}

func TestWorkflowRPCNodeTypeExclusivity(t *testing.T) {
	server, resourceType := newWorkflowRPCServer(t, nil)
	for _, tc := range []struct {
		name      string
		modify    func(configObject)
		unknownAt *tftypes.AttributePath
		want      []string
	}{
		{
			name:   "missing",
			modify: func(c configObject) { delete(rpcNode(c, 1), "webhook") },
			want:   []string{"Invalid node type combination"},
		},
		{
			name:   "multiple",
			modify: func(c configObject) { rpcNode(c, 1)["kafka"] = configObject{} },
			want:   []string{"Invalid node type combination"},
		},
		{
			name:      "unknown_block_defers",
			modify:    func(configObject) {},
			unknownAt: tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(1).WithAttributeName("condition"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := rpcWorkflowConfig(2)
			tc.modify(config)
			value := configValue(t, resourceType, config)
			if tc.unknownAt != nil {
				value = withUnknownValue(t, value, tc.unknownAt)
			}
			assert.Equal(t, tc.want, rpcErrorSummaries(validateResourceConfig(t, server, value)))
		})
	}
}

func TestWorkflowRPCValidationStillRejectsInvalidGraphs(t *testing.T) {
	server, resourceType := newWorkflowRPCServer(t, nil)
	for _, tc := range []struct {
		name   string
		modify func(configObject)
	}{
		{"long_identifier", func(c configObject) { c["identifier"] = strings.Repeat("x", 61) }},
		{"duplicate_node", func(c configObject) { rpcNode(c, 1)["identifier"] = "trigger" }},
		{"invalid_connection", func(c configObject) {
			c["connections"] = configList{configObject{"source_identifier": "trigger", "target_identifier": "absent"}}
		}},
		{"invalid_dataset_combinator", func(c configObject) {
			serviceUserProperty(c)["dataset"].(configObject)["combinator"] = "xor"
		}},
		{"conflicting_dataset_selectors", func(c configObject) {
			rules := serviceUserProperty(c)["dataset"].(configObject)["rules"].(configList)
			rules[0].(configObject)["relation"] = "owner"
		}},
		{"false_required", func(c configObject) { serviceUserProperty(c)["required"] = false }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := rpcWorkflowConfig(2)
			tc.modify(config)
			diagnostics := validateResourceConfig(t, server, configValue(t, resourceType, config))
			assert.NotEmpty(t, rpcErrorSummaries(diagnostics), "invalid graph passed full RPC validation")
		})
	}
}

func TestWorkflowRPCReadRoundTripAndNoOpPlan(t *testing.T) {
	ctx := context.Background()
	var workflow *cli.Workflow
	gets := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/workflows/rpc_workflow" {
			t.Errorf("unexpected API request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		gets++
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"ok": true, "workflow": workflow}))
	}))
	defer api.Close()

	server, resourceType := newWorkflowRPCServer(t, &cli.PortClient{Client: resty.New().SetBaseURL(api.URL), JSONEscapeHTML: true})
	config := configValue(t, resourceType, rpcWorkflowConfig(2))

	created, err := planResourceChange(t, server, config, tftypes.NewValue(resourceType, nil)).PlannedState.Unmarshal(resourceType)
	require.NoError(t, err)
	created, err = tftypes.Transform(created, func(valuePath *tftypes.AttributePath, value tftypes.Value) (tftypes.Value, error) {
		if valuePath.Equal(tftypes.NewAttributePath().WithAttributeName("id")) {
			return tftypes.NewValue(value.Type(), "rpc_workflow"), nil
		}
		return value, nil
	})
	require.NoError(t, err)

	schemaResp := &resource.SchemaResponse{}
	(&WorkflowResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResp)
	var model WorkflowModel
	diags := tfsdk.State{Schema: schemaResp.Schema, Raw: created}.Get(ctx, &model)
	require.False(t, diags.HasError(), diags.Errors())
	workflow, err = workflowStateToPortBody(ctx, &model)
	require.NoError(t, err)

	read, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{TypeName: "port_workflow", CurrentState: dynamicValue(t, created)})
	require.NoError(t, err)
	requireNoErrorDiagnostics(t, read.Diagnostics)
	refreshed, err := read.NewState.Unmarshal(resourceType)
	require.NoError(t, err)
	require.Equal(t, 1, gets, "refresh must make exactly one API GET")

	replanned, err := planResourceChange(t, server, config, refreshed).PlannedState.Unmarshal(resourceType)
	require.NoError(t, err)
	assert.True(t, replanned.Equal(refreshed), "unchanged refreshed workflow did not produce a no-op plan")
}

// BenchmarkWorkflowRPCPlan runs both RPCs Terraform sends for a resource
// during `terraform plan`.
func BenchmarkWorkflowRPCPlan(b *testing.B) {
	for _, nodes := range []int{2, 16, 32} {
		b.Run(fmt.Sprintf("nodes_%d", nodes), func(b *testing.B) {
			server, resourceType := newWorkflowRPCServer(b, nil)
			config := configValue(b, resourceType, rpcWorkflowConfig(nodes))
			prior := tftypes.NewValue(resourceType, nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				requireNoErrorDiagnostics(b, validateResourceConfig(b, server, config))
				planResourceChange(b, server, config, prior)
			}
		})
	}
}
