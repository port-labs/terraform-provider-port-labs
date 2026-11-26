package cli

// scorecardRequestsAdditionalProperties reports whether the Port API must persist
// custom _scorecard / _rule entity fields for this scorecard payload.
func scorecardRequestsAdditionalProperties(scorecard *Scorecard) bool {
	if scorecard == nil {
		return false
	}
	if len(scorecard.Properties) > 0 || len(scorecard.Relations) > 0 {
		return true
	}
	for _, rule := range scorecard.Rules {
		if len(rule.Properties) > 0 || len(rule.Relations) > 0 {
			return true
		}
	}
	return false
}