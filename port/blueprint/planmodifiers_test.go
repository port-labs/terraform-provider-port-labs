package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// planComputed runs the plan modifiers of the computed string attribute attrName against a prior state and a plan in
// which only the given string attributes (identifier, blueprint, ...) are set, and returns the planned value.
func planComputed(t *testing.T, attr schema.StringAttribute, attrName string, prior string, priorAttrs, plannedAttrs map[string]string) types.String {
	t.Helper()
	ctx := context.Background()
	attrs := map[string]schema.Attribute{attrName: schema.StringAttribute{Computed: true}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{attrName: tftypes.String}}
	for k := range priorAttrs {
		attrs[k] = schema.StringAttribute{Optional: true}
		objType.AttributeTypes[k] = tftypes.String
	}
	s := schema.Schema{Attributes: attrs}
	build := func(self tftypes.Value, vals map[string]string) tftypes.Value {
		m := map[string]tftypes.Value{attrName: self}
		for k, v := range vals {
			m[k] = tftypes.NewValue(tftypes.String, v)
		}
		return tftypes.NewValue(objType, m)
	}
	req := planmodifier.StringRequest{
		Path:        path.Root(attrName),
		ConfigValue: types.StringNull(),
		PlanValue:   types.StringUnknown(),
		Plan:        tfsdk.Plan{Schema: s, Raw: build(tftypes.NewValue(tftypes.String, tftypes.UnknownValue), plannedAttrs)},
		StateValue:  types.StringValue(prior),
		State:       tfsdk.State{Schema: s, Raw: build(tftypes.NewValue(tftypes.String, prior), priorAttrs)},
	}
	resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
	for _, m := range attr.PlanModifiers {
		m.PlanModifyString(ctx, req, resp)
		req.PlanValue = resp.PlanValue
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("plan modifier error: %v", resp.Diagnostics)
	}
	return resp.PlanValue
}

func TestBlueprintIdPlan(t *testing.T) {
	id := BlueprintSchema()["id"].(schema.StringAttribute)
	same := map[string]string{"identifier": "svc"}

	if got := planComputed(t, id, "id", "svc", same, same); got != types.StringValue("svc") {
		t.Errorf("unchanged identifier: want prior id kept, got %v", got)
	}
	// TestAccPortUpdateBlueprintIdentifier / TestAccPortWebhookCreateWithInvalidRelations rename the blueprint.
	if got := planComputed(t, id, "id", "svc", same, map[string]string{"identifier": "svc2"}); !got.IsUnknown() {
		t.Errorf("renamed identifier: want unknown id, got %v", got)
	}
}
