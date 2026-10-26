package webhook

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshWebhookStateKeepsSecretWhenAPIOmitsIt(t *testing.T) {
	now := time.Date(2026, 10, 26, 9, 0, 0, 0, time.UTC)
	header := "X-Hub-Signature-256"
	algorithm := "plain"
	state := &WebhookModel{
		Security: &SecurityModel{
			Secret:              types.StringValue("configured-secret"),
			SignatureHeaderName: types.StringValue("old-header"),
			SignatureAlgorithm:  types.StringValue("sha256"),
		},
	}
	api := &cli.Webhook{
		Meta: cli.Meta{
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: "creator",
			UpdatedBy: "updater",
		},
		Identifier: "pr-webhook",
		Url:        "https://example.com/hook",
		WebhookKey: "whk",
		Security: &cli.Security{
			SignatureHeaderName: &header,
			SignatureAlgorithm:  &algorithm,
		},
	}

	if err := refreshWebhookState(context.Background(), state, api); err != nil {
		t.Fatalf("refreshWebhookState: %v", err)
	}

	if state.Security == nil {
		t.Fatal("security block was cleared")
	}
	if got := state.Security.Secret.ValueString(); got != "configured-secret" {
		t.Errorf("secret = %q, want configured-secret", got)
	}
	if got := state.Security.SignatureHeaderName.ValueString(); got != header {
		t.Errorf("signature_header_name = %q, want %q", got, header)
	}
	if got := state.Security.SignatureAlgorithm.ValueString(); got != algorithm {
		t.Errorf("signature_algorithm = %q, want %q", got, algorithm)
	}
}

func TestWebhookStableComputedUseStateForUnknown(t *testing.T) {
	attrs := WebhookSchema()
	for _, name := range []string{"id", "url", "webhook_key", "created_at", "created_by", "updated_at", "updated_by"} {
		attr, ok := attrs[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", name)
		}
		if !usesStateForUnknown(attr.PlanModifiers) {
			t.Errorf("%s is missing UseStateForUnknown", name)
		}
	}
}

func usesStateForUnknown(modifiers []planmodifier.String) bool {
	ctx := context.Background()
	for _, modifier := range modifiers {
		resp := &planmodifier.StringResponse{PlanValue: types.StringUnknown()}
		modifier.PlanModifyString(ctx, planmodifier.StringRequest{
			StateValue: types.StringValue("keep"),
			PlanValue:  types.StringUnknown(),
		}, resp)
		if resp.PlanValue.ValueString() == "keep" {
			return true
		}
	}
	return false
}
