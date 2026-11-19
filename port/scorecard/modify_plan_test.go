package scorecard

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func testModifyPlanRule(id string) Rule {
	return Rule{
		Identifier: types.StringValue(id),
		Title:      types.StringValue(id + "-title"),
		Level:      types.StringValue("Gold"),
	}
}

func TestIsRulesOrderOnlyChange(t *testing.T) {
	t.Parallel()

	a := []Rule{testModifyPlanRule("zebra"), testModifyPlanRule("alpha"), testModifyPlanRule("beta")}
	b := []Rule{testModifyPlanRule("beta"), testModifyPlanRule("zebra"), testModifyPlanRule("alpha")}
	assert.True(t, isRulesOrderOnlyChange(a, b))
	assert.False(t, isRulesOrderOnlyChange(a, a))

	b[0].Title = types.StringValue("changed")
	assert.False(t, isRulesOrderOnlyChange(a, b))

	assert.False(t, isRulesOrderOnlyChange(a, a[:2]))

	differentIDs := []Rule{testModifyPlanRule("beta"), testModifyPlanRule("zebra"), testModifyPlanRule("gamma")}
	assert.False(t, isRulesOrderOnlyChange(a, differentIDs))
}
