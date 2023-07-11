package flex

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoStringToFramework(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}

	return types.StringValue(*v)
}
