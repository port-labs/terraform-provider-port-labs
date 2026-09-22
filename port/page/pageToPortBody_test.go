package page

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPageToPortBodyPreservesIframeWidgetRootRelativeURL(t *testing.T) {
	relativeURL := "/~/apiEntity?identifier=ads-api"
	widgetJSON := `{"type":"dashboard-widget","id":"dashboardWidget","widgets":[{"type":"iframe-widget","id":"embeddedView","url":"` + relativeURL + `","urlType":"public"}]}`

	state := &PageModel{
		Identifier: types.StringValue("embedded_port_page"),
		Type:       types.StringValue("dashboard"),
		Widgets: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(widgetJSON),
		}),
	}

	body, err := PageToPortBody(state)
	require.NoError(t, err)
	require.NotNil(t, body.Widgets)
	require.Len(t, *body.Widgets, 1)

	dashboardWidget := (*body.Widgets)[0]
	nestedWidgets, ok := dashboardWidget["widgets"].([]any)
	require.True(t, ok)
	require.Len(t, nestedWidgets, 1)

	iframeWidget, ok := nestedWidgets[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, relativeURL, iframeWidget["url"])
}
