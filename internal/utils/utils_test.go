package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoObjectToTerraformString(t *testing.T) {
	type args struct {
		v              interface{}
		JSONEscapeHTML bool
	}
	tests := []struct {
		name    string
		args    args
		want    types.String
		wantErr bool
	}{
		{name: "escape-html nil map", args: args{v: (map[any]any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil *map", args: args{v: (*map[any]any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil array", args: args{v: (*[0]any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil chan", args: args{v: (chan any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil *chan", args: args{v: (*chan any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil slice", args: args{v: ([]any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil *slice", args: args{v: (*[]any)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{name: "escape-html nil *string", args: args{v: (*string)(nil), JSONEscapeHTML: true}, want: types.StringNull()},
		{
			name: "escape-html nested map",
			args: args{v: map[string]any{
				"string": "hello world",
				"int":    42,
				"bool":   true,
				"float":  42.42,
				"list":   []any{1, 2, 3},
				"nil":    nil,
				"nested": map[string]any{"foo": "bar"},
			}, JSONEscapeHTML: true},
			want: types.StringValue(
				"{\"bool\":true,\"float\":42.42,\"int\":42,\"list\":[1,2,3],\"nested\":{\"foo\":\"bar\"},\"nil\":null,\"string\":\"hello world\"}",
			),
		},
		{
			name: "escape-html nested slice",
			args: args{v: []any{"hello world", 1, map[string]any{"foo": "bar"}}, JSONEscapeHTML: true},
			want: types.StringValue("[\"hello world\",1,{\"foo\":\"bar\"}]"),
		},
		{
			name: "escape-html html unsafe characters",
			args: args{v: map[string]any{"property": "sum", "operator": ">", "value": 2}, JSONEscapeHTML: true},
			want: types.StringValue("{\"operator\":\"\\u003e\",\"property\":\"sum\",\"value\":2}"),
		},
		{
			name: "escape-html html unsafe characters (double escape)",
			args: args{v: "{\"property\": \"sum\", \"operator\": \">\", \"value\": 2}", JSONEscapeHTML: true},
			want: types.StringValue("\"{\\\"property\\\": \\\"sum\\\", \\\"operator\\\": \\\"\\u003e\\\", \\\"value\\\": 2}\""),
		},
		{name: "no-escape-html nil map", args: args{v: (map[any]any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil *map", args: args{v: (*map[any]any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil array", args: args{v: (*[0]any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil chan", args: args{v: (chan any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil *chan", args: args{v: (*chan any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil slice", args: args{v: ([]any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil *slice", args: args{v: (*[]any)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{name: "no-escape-html nil *string", args: args{v: (*string)(nil), JSONEscapeHTML: false}, want: types.StringNull()},
		{
			name: "no-escape-html nested map",
			args: args{v: map[string]any{
				"string": "hello world",
				"int":    42,
				"bool":   true,
				"float":  42.42,
				"list":   []any{1, 2, 3},
				"nil":    nil,
				"nested": map[string]any{"foo": "bar"},
			}, JSONEscapeHTML: false},
			want: types.StringValue(
				"{\"bool\":true,\"float\":42.42,\"int\":42,\"list\":[1,2,3],\"nested\":{\"foo\":\"bar\"},\"nil\":null,\"string\":\"hello world\"}",
			),
		},
		{
			name: "no-escape-html nested slice",
			args: args{v: []any{"hello world", 1, map[string]any{"foo": "bar"}}, JSONEscapeHTML: false},
			want: types.StringValue("[\"hello world\",1,{\"foo\":\"bar\"}]"),
		},
		{
			name: "no-escape-html html unsafe characters",
			args: args{v: map[string]any{"property": "sum", "operator": ">", "value": 2}, JSONEscapeHTML: false},
			want: types.StringValue("{\"operator\":\">\",\"property\":\"sum\",\"value\":2}"),
		},
		{
			name: "no-escape-html html unsafe characters (double escape)",
			args: args{v: "{\"property\": \"sum\", \"operator\": \">\", \"value\": 2}", JSONEscapeHTML: false},
			want: types.StringValue("\"{\\\"property\\\": \\\"sum\\\", \\\"operator\\\": \\\">\\\", \\\"value\\\": 2}\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoObjectToTerraformString(tt.args.v, tt.args.JSONEscapeHTML)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
