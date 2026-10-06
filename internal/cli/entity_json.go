package cli

import "encoding/json"

func (e *Entity) UnmarshalJSON(data []byte) error {
	type entityAlias Entity
	if err := json.Unmarshal(data, (*entityAlias)(e)); err != nil {
		return err
	}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if e.RelationSources == nil {
		e.RelationSources = unmarshalEntitySourceMap(raw, "relationSources", "_relationSources")
	}
	if e.TeamSources == nil {
		e.TeamSources = unmarshalEntitySourceMap(raw, "teamSources", "_teamSources")
	}

	return nil
}

func unmarshalEntitySourceMap(raw map[string]json.RawMessage, keys ...string) map[string]any {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}

		var sources map[string]any
		if err := json.Unmarshal(value, &sources); err != nil {
			continue
		}

		return sources
	}

	return nil
}

// EntityTeamIdentifiers returns flattened team identifiers when the API returns a string array.
func EntityTeamIdentifiers(team any) []string {
	switch value := team.(type) {
	case nil:
		return nil
	case []string:
		return value
	case []any:
		teams := make([]string, 0, len(value))
		for _, item := range value {
			if teamID, ok := item.(string); ok {
				teams = append(teams, teamID)
			}
		}
		return teams
	default:
		return nil
	}
}

func EntityTeamIsUnionSlice(team any) bool {
	_, ok := team.(map[string]any)
	return ok
}
