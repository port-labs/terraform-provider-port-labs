package entity

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func blueprintRelationIsUnion(blueprint *cli.Blueprint, relationIdentifier string) bool {
	if blueprint == nil {
		return false
	}

	relation, ok := blueprint.Relations[relationIdentifier]
	if !ok {
		return false
	}

	return relation.Union != nil && *relation.Union
}

func blueprintHasUnionTeamOwnership(blueprint *cli.Blueprint) bool {
	if blueprint == nil || blueprint.Ownership == nil {
		return false
	}

	if blueprint.Ownership.Type != "Direct" {
		return false
	}

	return blueprint.Ownership.Union != nil && *blueprint.Ownership.Union
}
