package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoObjectToTerraformString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    types.String
		wantErr bool
	}{
		{name: "nil map", args: args{v: (map[any]any)(nil)}, want: types.StringNull()},
		{name: "nil *map", args: args{v: (*map[any]any)(nil)}, want: types.StringNull()},
		{name: "nil array", args: args{v: (*[0]any)(nil)}, want: types.StringNull()},
		{name: "nil chan", args: args{v: (chan any)(nil)}, want: types.StringNull()},
		{name: "nil *chan", args: args{v: (*chan any)(nil)}, want: types.StringNull()},
		{name: "nil slice", args: args{v: ([]any)(nil)}, want: types.StringNull()},
		{name: "nil *slice", args: args{v: (*[]any)(nil)}, want: types.StringNull()},
		{name: "nil *string", args: args{v: (*string)(nil)}, want: types.StringNull()},

		{
			name: "nested map",
			args: args{v: map[string]any{
				"string": "hello world",
				"int":    42,
				"bool":   true,
				"float":  42.42,
				"list":   []any{1, 2, 3},
				"nil":    nil,
				"nested": map[string]any{"foo": "bar"},
			}},
			want: types.StringValue(
				"{\"bool\":true,\"float\":42.42,\"int\":42,\"list\":[1,2,3],\"nested\":{\"foo\":\"bar\"},\"nil\":null,\"string\":\"hello world\"}\n",
			),
		},

		{
			name: "nested slice",
			args: args{v: []any{"hello world", 1, map[string]any{"foo": "bar"}}},
			want: types.StringValue("[\"hello world\",1,{\"foo\":\"bar\"}]\n"),
		},

		{
			name: "html unsafe characters",
			args: args{v: map[string]any{"property": "sum", "operator": ">", "value": 2}},
			want: types.StringValue("{\"operator\":\">\",\"property\":\"sum\",\"value\":2}\n"),
		},
		{
			name: "html unsafe characters (double escape)",
			args: args{v: "{\"property\": \"sum\", \"operator\": \">\", \"value\": 2}"},
			want: types.StringValue("\"{\\\"property\\\": \\\"sum\\\", \\\"operator\\\": \\\">\\\", \\\"value\\\": 2}\"\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoObjectToTerraformString(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
