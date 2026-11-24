package virtual_mcp_tool

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/require"
)

func TestVirtualMcpToolResourceToPortBodyListedArguments(t *testing.T) {
	state := &VirtualMcpToolModel{
		Title:           types.StringValue("list_skills"),
		Name:            types.StringValue("list_skills"),
		Tool:            types.StringValue("list_entities"),
		ExposeArguments: types.StringValue(ExposeArgumentsListed),
		Arguments:       types.StringValue(`{"blueprintIdentifier":{"value":"skill"}}`),
	}

	entity, err := virtualMcpToolResourceToPortBody(state)
	require.NoError(t, err)
	require.Equal(t, VirtualMcpToolBlueprintIdentifier, entity.Blueprint)
	require.Equal(t, ExposeArgumentsListed, entity.Properties["expose_arguments"])
	arguments, ok := entity.Properties["arguments"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, arguments["blueprintIdentifier"])
}

func TestRefreshVirtualMcpToolStateDefaultsExposeArguments(t *testing.T) {
	entity := &cli.Entity{
		Identifier: "skills_list",
		Title:      "list_skills",
		Properties: map[string]any{
			"name": "list_skills",
			"tool": "list_entities",
		},
	}

	state := &VirtualMcpToolModel{}
	err := refreshVirtualMcpToolState(context.Background(), state, entity)
	require.NoError(t, err)
	require.Equal(t, ExposeArgumentsAll, state.ExposeArguments.ValueString())
}