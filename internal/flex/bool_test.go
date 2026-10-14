package flex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoBoolToFrameworkDefaultFalse(t *testing.T) {
	assert.True(t, GoBoolToFrameworkDefaultFalse(nil).Equal(types.BoolValue(false)))
	assert.True(t, GoBoolToFrameworkDefaultFalse(ptr(true)).Equal(types.BoolValue(true)))

	falseValue := false
	assert.True(t, GoBoolToFrameworkDefaultFalse(&falseValue).Equal(types.BoolValue(false)))
}

func TestFrameworkBoolToTruePointer(t *testing.T) {
	assert.Nil(t, FrameworkBoolToTruePointer(types.BoolNull()))
	assert.Nil(t, FrameworkBoolToTruePointer(types.BoolValue(false)))

	truePointer := FrameworkBoolToTruePointer(types.BoolValue(true))
	require.NotNil(t, truePointer)
	assert.True(t, *truePointer)
}

func ptr(v bool) *bool {
	return &v
}
