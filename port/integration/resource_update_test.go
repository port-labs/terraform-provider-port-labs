package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateOmitsUnchangedSpecAndSkipsInfraWait(t *testing.T) {
	elementType := notificationChannelElementType()
	sameSpec := types.StringValue(`{"integrationSpec":{"token":"secret"}}`)
	rotatedSpec := types.StringValue(`{"integrationSpec":{"token":"rotated"}}`)
	emailOn := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: true})
	emailOff := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: false})

	tests := []struct {
		name            string
		state           *IntegrationModel
		plan            *IntegrationModel
		wantSpec        bool
		wantIntegration int
		wantOptInPost   int
	}{
		{
			name:            "title change keeps spec off the body",
			state:           hostedIntegrationModel("before", types.MapNull(elementType)),
			plan:            hostedIntegrationModel("after", types.MapNull(elementType)),
			wantIntegration: 0,
		},
		{
			name:            "notification-only change keeps spec off the body",
			state:           hostedIntegrationModel("title", emailOn),
			plan:            hostedIntegrationModel("title", emailOff),
			wantIntegration: 0,
			wantOptInPost:   1,
		},
		{
			name:            "spec change is sent and waited on",
			state:           withSpec(hostedIntegrationModel("title", types.MapNull(elementType)), sameSpec),
			plan:            withSpec(hostedIntegrationModel("title", types.MapNull(elementType)), rotatedSpec),
			wantSpec:        true,
			wantIntegration: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var patchBody map[string]any
			var integrationGets, optInPosts int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if strings.Contains(r.URL.Path, "/notifications/opt-ins") {
					if r.Method == http.MethodPost {
						optInPosts++
						_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"email","enabled":false}}`))
						return
					}
					_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"email","enabled":true,"emails":[],"teams":[]}]}`))
					return
				}
				if r.Method == http.MethodGet {
					integrationGets++
				}
				if r.Method == http.MethodPatch {
					raw, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.NoError(t, json.Unmarshal(raw, &patchBody))
				}
				_, _ = w.Write([]byte(runningIntegrationJSON))
			}))
			t.Cleanup(server.Close)

			integration := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
			diags := updateIntegrationModels(t, integration, tt.state, tt.plan)
			require.False(t, diags.HasError(), diags)

			_, hasSpec := patchBody["spec"]
			assert.Equal(t, tt.wantSpec, hasSpec, "spec present in update body: %v", patchBody["spec"])
			assert.Equal(t, tt.wantIntegration, integrationGets, "integration reads while waiting")
			assert.Equal(t, tt.wantOptInPost, optInPosts, "opt-in posts")
		})
	}
}

func withSpec(model *IntegrationModel, spec types.String) *IntegrationModel {
	model.Spec = spec
	return model
}

func updateIntegrationModels(t *testing.T, integration *IntegrationResource, state, plan *IntegrationModel) diag.Diagnostics {
	t.Helper()
	resourceSchema := integrationResourceSchema(t)
	req := resource.UpdateRequest{
		State: tfsdk.State{Schema: resourceSchema},
		Plan:  tfsdk.Plan{Schema: resourceSchema},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema}}
	diags := req.State.Set(context.Background(), state)
	diags.Append(req.Plan.Set(context.Background(), plan)...)
	require.False(t, diags.HasError(), diags)
	integration.Update(context.Background(), req, resp)
	return resp.Diagnostics
}
