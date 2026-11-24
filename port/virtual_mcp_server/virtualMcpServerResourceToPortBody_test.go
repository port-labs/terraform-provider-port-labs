package virtual_mcp_server

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestVirtualMcpServerResourceToPortBody(t *testing.T) {
	state := &VirtualMcpServerModel{
		Title:            types.StringValue("Skills Distribution"),
		Description:      types.StringValue("Distributes skills"),
		Instructions:     types.StringValue("Call list_skills first."),
		PassthroughTools: []types.String{types.StringValue("load_skill")},
		VirtualMcpTools:  []types.String{types.StringValue("skills_distribution_list_skills")},
	}

	entity, err := virtualMcpServerResourceToPortBody(state)
	require.NoError(t, err)
	require.Equal(t, VirtualMcpServerBlueprintIdentifier, entity.Blueprint)
	require.Equal(t, []string{"load_skill"}, entity.Properties["passthrough_tools"])
	require.Equal(t, []string{"skills_distribution_list_skills"}, entity.Relations[VirtualMcpToolsRelation])
}

func TestRefreshVirtualMcpServerStateReadsURL(t *testing.T) {
	entity := &cli.Entity{
		Identifier: "skills_distribution",
		Title:      "Skills Distribution",
		Properties: map[string]any{
			"url": "https://mcp.port.io/v1/virtual-server/skills_distribution",
		},
		Relations: map[string]any{
			VirtualMcpToolsRelation: []any{"skills_distribution_list_skills"},
		},
	}

	state := &VirtualMcpServerModel{}
	err := refreshVirtualMcpServerState(context.Background(), state, entity)
	require.NoError(t, err)
	require.Equal(t, "https://mcp.port.io/v1/virtual-server/skills_distribution", state.URL.ValueString())
	require.Equal(t, []string{"skills_distribution_list_skills"}, toolIDs(state))
}

func toolIDs(state *VirtualMcpServerModel) []string {
	out := make([]string, len(state.VirtualMcpTools))
	for i, id := range state.VirtualMcpTools {
		out[i] = id.ValueString()
	}
	return out
}