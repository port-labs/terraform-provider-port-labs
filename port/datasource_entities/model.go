package datasource_entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EntityModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Blueprint  types.String `tfsdk:"blueprint"`
}

type DatasourceEntitiesDataModel struct {
	ID               types.String   `tfsdk:"id"`
	DatasourcePrefix types.String   `tfsdk:"datasource_prefix"`
	DatasourceSuffix types.String   `tfsdk:"datasource_suffix"`
	Limit            types.Int64    `tfsdk:"limit"`
	Before           types.String   `tfsdk:"before"`
	Entities         []EntityModel  `tfsdk:"entities"`
}

func (m *DatasourceEntitiesDataModel) GenerateID() string {
	var sb strings.Builder
	sb.WriteString(m.DatasourcePrefix.ValueString())
	sb.WriteString(m.DatasourceSuffix.ValueString())
	if !m.Limit.IsNull() {
		sb.WriteString(fmt.Sprintf("%d", m.Limit.ValueInt64()))
	}
	if !m.Before.IsNull() {
		sb.WriteString(m.Before.ValueString())
	}

	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}
