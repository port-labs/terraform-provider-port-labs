package scorecard

const (
	scorecardManagedRelationKey = "group"
	ruleManagedRelationKey      = "scorecard"
)

func omitManagedRelationKeys(relations map[string]any, keys ...string) map[string]any {
	if len(relations) == 0 {
		return relations
	}
	sanitized := make(map[string]any, len(relations))
	for key, value := range relations {
		skip := false
		for _, reserved := range keys {
			if key == reserved {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		sanitized[key] = value
	}
	return sanitized
}