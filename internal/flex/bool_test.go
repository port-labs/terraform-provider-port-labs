package flex

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGoBoolToFrameworkDefaultFalse(t *testing.T) {
	assert.True(t, GoBoolToFrameworkDefaultFalse(nil).Equal(types.BoolValue(false)))
	assert.True(t, GoBoolToFrameworkDefaultFalse(ptr(true)).Equal(types.BoolValue(true)))

	falseValue := false
	assert.True(t, GoBoolToFrameworkDefaultFalse(&falseValue).Equal(types.BoolValue(false)))
}

func ptr(v bool) *bool {
	return &v
}
