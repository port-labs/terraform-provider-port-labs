package scorecard_group

// scorecardGroupMembershipChanged reports whether blueprint membership changed between
// states. Membership changes require PUT because PATCH does not accept `blueprints` and
// only patches existing per-blueprint `scorecards` members.
func scorecardGroupMembershipChanged(previous, current *ScorecardGroupModel) bool {
	prevHasScorecards := len(previous.Scorecards) > 0
	currHasScorecards := len(current.Scorecards) > 0

	if prevHasScorecards != currHasScorecards {
		return true
	}

	if currHasScorecards {
		if len(previous.Scorecards) != len(current.Scorecards) {
			return true
		}
		for blueprintID := range current.Scorecards {
			if _, ok := previous.Scorecards[blueprintID]; !ok {
				return true
			}
		}
		for blueprintID := range previous.Scorecards {
			if _, ok := current.Scorecards[blueprintID]; !ok {
				return true
			}
		}
		return false
	}

	if len(previous.Blueprints) != len(current.Blueprints) {
		return true
	}

	prevBlueprints := make(map[string]struct{}, len(previous.Blueprints))
	for _, blueprint := range previous.Blueprints {
		prevBlueprints[blueprint.ValueString()] = struct{}{}
	}
	for _, blueprint := range current.Blueprints {
		if _, ok := prevBlueprints[blueprint.ValueString()]; !ok {
			return true
		}
	}
	return false
}
