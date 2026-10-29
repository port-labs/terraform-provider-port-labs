package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func encryptionToBody(encryption types.String, clientSide *ClientSideEncryptionModel) any {
	if !encryption.IsNull() {
		return encryption.ValueString()
	}
	if clientSide != nil {
		return map[string]string{
			"algorithm": clientSide.Algorithm.ValueString(),
			"key":       clientSide.Key.ValueString(),
		}
	}
	return nil
}

func encryptionToState(value any) (types.String, *ClientSideEncryptionModel) {
	if value == nil {
		return types.StringNull(), nil
	}

	switch enc := value.(type) {
	case string:
		return types.StringValue(enc), nil
	case map[string]any:
		algorithm, _ := enc["algorithm"].(string)
		key, _ := enc["key"].(string)
		if algorithm != "" && key != "" {
			return types.StringNull(), &ClientSideEncryptionModel{
				Algorithm: types.StringValue(algorithm),
				Key:       types.StringValue(key),
			}
		}
	case map[string]string:
		if algorithm, ok := enc["algorithm"]; ok {
			if key, ok := enc["key"]; ok && algorithm != "" && key != "" {
				return types.StringNull(), &ClientSideEncryptionModel{
					Algorithm: types.StringValue(algorithm),
					Key:       types.StringValue(key),
				}
			}
		}
	}

	return types.StringNull(), nil
}
