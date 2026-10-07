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
)

type rpcProvider struct{ client *cli.PortClient }

func (p *rpcProvider) Metadata(_ context.Context, _ provider.MetadataRequest, r *provider.MetadataResponse) {
	r.TypeName = "port"
}
func (p *rpcProvider) Schema(_ context.Context, _ provider.SchemaRequest, r *provider.SchemaResponse) {
	r.Schema = providerschema.Schema{}
}
func (p *rpcProvider) Configure(_ context.Context, _ provider.ConfigureRequest, r *provider.ConfigureResponse) {
	r.ResourceData = p.client
}
func (p *rpcProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewWorkflowResource}
}
func (p *rpcProvider) DataSources(context.Context) []func() datasource.DataSource { return nil }

func rpcServer(t testing.TB, client *cli.PortClient) (tfprotov6.ProviderServer, tftypes.Type) {
	t.Helper()
	ctx := context.Background()
	s := providerserver.NewProtocol6(&rpcProvider{client: client})()
	schemas, err := s.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	rpcCheck(t, schemas.Diagnostics)
	c, err := tfprotov6.NewDynamicValue(schemas.Provider.ValueType(), tftypes.NewValue(schemas.Provider.ValueType(), map[string]tftypes.Value{}))
	if err != nil {
		t.Fatal(err)
	}
	configured, err := s.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{Config: &c})
	if err != nil {
		t.Fatal(err)
	}
	rpcCheck(t, configured.Diagnostics)
	return s, schemas.ResourceSchemas["port_workflow"].ValueType()
}
func rpcCheck(t testing.TB, ds []*tfprotov6.Diagnostic) {
	t.Helper()
	for _, d := range ds {
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("%s: %s", d.Summary, d.Detail)
		}
	}
}
func rpcDynamic(t testing.TB, v tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	d, e := tfprotov6.NewDynamicValue(v.Type(), v)
	if e != nil {
		t.Fatal(e)
	}
	return &d
}
func rpcValue(t testing.TB, typ tftypes.Type, data any) tftypes.Value {
	t.Helper()
	if data == nil {
		return tftypes.NewValue(typ, nil)
	}
	if data == tftypes.UnknownValue {
		return tftypes.NewValue(typ, data)
	}
	switch typ := typ.(type) {
	case tftypes.Object:
		values := map[string]tftypes.Value{}
		m := data.(map[string]any)
		for k, v := range typ.AttributeTypes {
			values[k] = rpcValue(t, v, m[k])
		}
		return tftypes.NewValue(typ, values)
	case tftypes.List:
		values := []tftypes.Value{}
		for _, v := range data.([]any) {
			values = append(values, rpcValue(t, typ.ElementType, v))
		}
		return tftypes.NewValue(typ, values)
	case tftypes.Map:
		values := map[string]tftypes.Value{}
		for k, v := range data.(map[string]any) {
			values[k] = rpcValue(t, typ.ElementType, v)
		}
		return tftypes.NewValue(typ, values)
	}
	return tftypes.NewValue(typ, data)
}
func rpcConfig(nodes int) map[string]any {
	n := []any{map[string]any{"identifier": "trigger", "self_serve_trigger": map[string]any{"user_inputs": map[string]any{"user_properties": map[string]any{"string_props": map[string]any{"service": map[string]any{"title": "Service", "format": "entity", "blueprint": "service", "dataset": map[string]any{"combinator": "and", "rules": []any{map[string]any{"property": "$identifier", "operator": "=", "value_json": "\"example\""}}}}}}}}}}
	for i := 1; i < nodes; i++ {
		n = append(n, map[string]any{"identifier": fmt.Sprintf("webhook_%d", i), "webhook": map[string]any{"url": "https://example.invalid/hook"}})
	}
	return map[string]any{"identifier": "rpc_workflow", "title": "RPC regression", "node": n}
}
func rpcAt(t testing.TB, v tftypes.Value, p *tftypes.AttributePath) tftypes.Value {
	t.Helper()
	x, remain, e := tftypes.WalkAttributePath(v, p)
	if e != nil {
		t.Fatalf("%s: %v (%s)", p, e, remain)
	}
	return x.(tftypes.Value)
}
func rpcPlan(t testing.TB, s tfprotov6.ProviderServer, typ tftypes.Type, config, prior tftypes.Value) *tfprotov6.PlanResourceChangeResponse {
	t.Helper()
	r, e := s.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{TypeName: "port_workflow", Config: rpcDynamic(t, config), ProposedNewState: rpcDynamic(t, config), PriorState: rpcDynamic(t, prior)})
	if e != nil {
		t.Fatal(e)
	}
	rpcCheck(t, r.Diagnostics)
	return r
}

func TestWorkflowRPCDefaultsNullsAndUnknowns(t *testing.T) {
	s, typ := rpcServer(t, nil)
	config := rpcValue(t, typ, rpcConfig(2))
	null := tftypes.NewValue(typ, nil)
	valid, e := s.ValidateResourceConfig(context.Background(), &tfprotov6.ValidateResourceConfigRequest{TypeName: "port_workflow", Config: rpcDynamic(t, config)})
	if e != nil {
		t.Fatal(e)
	}
	rpcCheck(t, valid.Diagnostics)
	r := rpcPlan(t, s, typ, config, null)
	plan, e := r.PlannedState.Unmarshal(typ)
	if e != nil {
		t.Fatal(e)
	}
	checks := []struct {
		p    *tftypes.AttributePath
		want any
	}{
		{tftypes.NewAttributePath().WithAttributeName("allow_anyone_to_view_runs"), true},
		{tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(0).WithAttributeName("verbose"), false},
		{tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(0).WithAttributeName("self_serve_trigger").WithAttributeName("published"), true},
		{tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(1).WithAttributeName("webhook").WithAttributeName("on_failure"), "terminate"},
	}
	for _, c := range checks {
		v := rpcAt(t, plan, c.p)
		if !v.Equal(tftypes.NewValue(v.Type(), c.want)) {
			t.Fatalf("default %s=%s", c.p, v)
		}
	}
	p := tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(0).WithAttributeName("self_serve_trigger").WithAttributeName("user_inputs").WithAttributeName("user_properties").WithAttributeName("number_props")
	if !rpcAt(t, plan, p).IsNull() {
		t.Fatal("absent numeric properties must remain null")
	}
	if rpcAt(t, plan, tftypes.NewAttributePath().WithAttributeName("id")).IsKnown() {
		t.Fatal("new resource ID must remain unknown")
	}
	for _, p := range []*tftypes.AttributePath{
		tftypes.NewAttributePath().WithAttributeName("node"),
		tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(1).WithAttributeName("webhook").WithAttributeName("url"),
		tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(0).WithAttributeName("condition"),
	} {
		unknown, e := tftypes.Transform(config, func(at *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if at.Equal(p) {
				return tftypes.NewValue(v.Type(), tftypes.UnknownValue), nil
			}
			return v, nil
		})
		if e != nil {
			t.Fatal(e)
		}
		valid, e := s.ValidateResourceConfig(context.Background(), &tfprotov6.ValidateResourceConfigRequest{TypeName: "port_workflow", Config: rpcDynamic(t, unknown)})
		if e != nil {
			t.Fatal(e)
		}
		rpcCheck(t, valid.Diagnostics)
	}
}

func TestWorkflowRPCValidationStillRejectsInvalidGraphs(t *testing.T) {
	s, typ := rpcServer(t, nil)
	for _, tc := range []struct {
		name   string
		modify func(map[string]any)
	}{
		{"long_identifier", func(c map[string]any) { c["identifier"] = strings.Repeat("x", 61) }},
		{"two_node_types", func(c map[string]any) {
			c["node"].([]any)[0].(map[string]any)["webhook"] = map[string]any{"url": "https://example.invalid"}
		}},
		{"duplicate_node", func(c map[string]any) { c["node"].([]any)[1].(map[string]any)["identifier"] = "trigger" }},
		{"invalid_connection", func(c map[string]any) {
			c["connections"] = []any{map[string]any{"source_identifier": "trigger", "target_identifier": "absent"}}
		}},
		{"invalid_dataset_combinator", func(c map[string]any) {
			service := c["node"].([]any)[0].(map[string]any)["self_serve_trigger"].(map[string]any)["user_inputs"].(map[string]any)["user_properties"].(map[string]any)["string_props"].(map[string]any)["service"].(map[string]any)
			service["dataset"].(map[string]any)["combinator"] = "xor"
		}},
		{"conflicting_dataset_selectors", func(c map[string]any) {
			service := c["node"].([]any)[0].(map[string]any)["self_serve_trigger"].(map[string]any)["user_inputs"].(map[string]any)["user_properties"].(map[string]any)["string_props"].(map[string]any)["service"].(map[string]any)
			service["dataset"].(map[string]any)["rules"].([]any)[0].(map[string]any)["relation"] = "owner"
		}},
		{"false_required", func(c map[string]any) {
			service := c["node"].([]any)[0].(map[string]any)["self_serve_trigger"].(map[string]any)["user_inputs"].(map[string]any)["user_properties"].(map[string]any)["string_props"].(map[string]any)["service"].(map[string]any)
			service["required"] = false
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := rpcConfig(2)
			tc.modify(c)
			r, e := s.ValidateResourceConfig(context.Background(), &tfprotov6.ValidateResourceConfigRequest{TypeName: "port_workflow", Config: rpcDynamic(t, rpcValue(t, typ, c))})
			if e != nil {
				t.Fatal(e)
			}
			for _, d := range r.Diagnostics {
				if d.Severity == tfprotov6.DiagnosticSeverityError {
					return
				}
			}
			t.Fatal("invalid graph passed full RPC validation")
		})
	}
}

func TestWorkflowRPCReadRoundTripAndNoOpPlan(t *testing.T) {
	ctx := context.Background()
	var body *cli.Workflow
	gets := 0
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/workflows/rpc_workflow" {
			t.Errorf("unexpected API request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(400)
			return
		}
		gets++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "workflow": body})
	}))
	defer httpServer.Close()
	s, typ := rpcServer(t, &cli.PortClient{Client: resty.New().SetBaseURL(httpServer.URL), JSONEscapeHTML: true})
	config := rpcValue(t, typ, rpcConfig(2))
	r := rpcPlan(t, s, typ, config, tftypes.NewValue(typ, nil))
	raw, e := r.PlannedState.Unmarshal(typ)
	if e != nil {
		t.Fatal(e)
	}
	raw, e = tftypes.Transform(raw, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if p.Equal(tftypes.NewAttributePath().WithAttributeName("id")) {
			return tftypes.NewValue(v.Type(), "rpc_workflow"), nil
		}
		return v, nil
	})
	if e != nil {
		t.Fatal(e)
	}
	schemaResponse := &resource.SchemaResponse{}
	(&WorkflowResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResponse)
	var model WorkflowModel
	state := tfsdk.State{Schema: schemaResponse.Schema, Raw: raw}
	if d := state.Get(ctx, &model); d.HasError() {
		t.Fatal(d)
	}
	body, e = workflowStateToPortBody(ctx, &model)
	if e != nil {
		t.Fatal(e)
	}
	read, e := s.ReadResource(ctx, &tfprotov6.ReadResourceRequest{TypeName: "port_workflow", CurrentState: rpcDynamic(t, raw)})
	if e != nil {
		t.Fatal(e)
	}
	rpcCheck(t, read.Diagnostics)
	refreshed, e := read.NewState.Unmarshal(typ)
	if e != nil {
		t.Fatal(e)
	}
	if gets != 1 {
		t.Fatalf("refresh must make exactly one API GET; got %d", gets)
	}
	next := rpcPlan(t, s, typ, config, refreshed)
	planned, e := next.PlannedState.Unmarshal(typ)
	if e != nil {
		t.Fatal(e)
	}
	if !planned.Equal(refreshed) {
		t.Fatal("unchanged refreshed workflow did not produce a no-op Plan")
	}
}

func BenchmarkWorkflowRPCPlan(b *testing.B) {
	for _, nodes := range []int{2, 16, 32} {
		b.Run(fmt.Sprintf("nodes_%d", nodes), func(b *testing.B) {
			s, typ := rpcServer(b, nil)
			c := rpcValue(b, typ, rpcConfig(nodes))
			prior := tftypes.NewValue(typ, nil)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rpcPlan(b, s, typ, c, prior)
			}
		})
	}
}
