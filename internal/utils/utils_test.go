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
		{
			name:    "encode error",
			args:    args{v: func() {}},
			want:    types.StringNull(),
			wantErr: true,
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

func TestJSONStringsSemanticallyEqual(t *testing.T) {
	t.Run("identical strings", func(t *testing.T) {
		s := `{"a":1,"b":2}`
		equal, err := JSONStringsSemanticallyEqual(s, s, false)
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("different key order", func(t *testing.T) {
		a := `{"resources":[],"deleteDependentEntities":true}`
		b := `{"deleteDependentEntities":true,"resources":[]}`
		equal, err := JSONStringsSemanticallyEqual(a, b, false)
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("nested key reordering", func(t *testing.T) {
		a := `{"outer":{"zebra":1,"alpha":2}}`
		b := `{"outer":{"alpha":2,"zebra":1}}`
		equal, err := JSONStringsSemanticallyEqual(a, b, false)
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("html escape equivalence", func(t *testing.T) {
		a := `{"operator":">","property":"sum","value":2}`
		b := `{"operator":"\u003e","property":"sum","value":2}`
		equal, err := JSONStringsSemanticallyEqual(a, b, true)
		assert.NoError(t, err)
		assert.True(t, equal)
	})

	t.Run("semantically different", func(t *testing.T) {
		a := `{"a":1}`
		b := `{"a":2}`
		equal, err := JSONStringsSemanticallyEqual(a, b, false)
		assert.NoError(t, err)
		assert.False(t, equal)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := JSONStringsSemanticallyEqual(`{invalid`, `{"a":1}`, false)
		assert.Error(t, err)
	})
}

func TestGoObjectToTerraformStringPreferExisting(t *testing.T) {
	v := map[string]any{
		"deleteDependentEntities": true,
		"resources":                 []any{},
	}

	t.Run("returns preferred when semantically equal", func(t *testing.T) {
		preferred := types.StringValue(`{"resources":[],"deleteDependentEntities":true}`)
		got, err := GoObjectToTerraformStringPreferExisting(preferred, v, false)
		assert.NoError(t, err)
		assert.Equal(t, preferred, got)
	})

	t.Run("returns encoded when preferred is null", func(t *testing.T) {
		got, err := GoObjectToTerraformStringPreferExisting(types.StringNull(), v, false)
		assert.NoError(t, err)
		want, _ := GoObjectToTerraformString(v, false)
		assert.Equal(t, want, got)
	})

	t.Run("returns encoded when semantically different", func(t *testing.T) {
		preferred := types.StringValue(`{"deleteDependentEntities":false,"resources":[]}`)
		got, err := GoObjectToTerraformStringPreferExisting(preferred, v, false)
		assert.NoError(t, err)
		want, _ := GoObjectToTerraformString(v, false)
		assert.Equal(t, want, got)
	})

	t.Run("falls back to encoded when preferred is invalid json", func(t *testing.T) {
		preferred := types.StringValue(`{invalid`)
		got, err := GoObjectToTerraformStringPreferExisting(preferred, v, false)
		assert.NoError(t, err)
		want, _ := GoObjectToTerraformString(v, false)
		assert.Equal(t, want, got)
	})
}
