package scorecard_group

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func TestShouldRefreshGroupLevels(t *testing.T) {
	t.Parallel()

	customLevels := []scorecard.Level{
		{Title: types.StringValue("Ready"), Color: types.StringValue("green")},
	}
	customCliLevels := []cli.Level{{Title: "Ready", Color: "green"}}

	tests := []struct {
		name        string
		stateLevels []scorecard.Level
		cliLevels   []cli.Level
		want        bool
	}{
		{
			name:        "omit default levels when state has none",
			stateLevels: nil,
			cliLevels:   scorecard.DefaultCliLevels(),
			want:        false,
		},
		{
			name:        "refresh custom levels from api",
			stateLevels: customLevels,
			cliLevels:   customCliLevels,
			want:        true,
		},
		{
			name:        "refresh non-default api levels when state has none",
			stateLevels: nil,
			cliLevels:   customCliLevels,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := shouldRefreshGroupLevels(tt.stateLevels, tt.cliLevels); got != tt.want {
				t.Fatalf("shouldRefreshGroupLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}
