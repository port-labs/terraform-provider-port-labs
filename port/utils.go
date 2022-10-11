package port

import "github.com/samber/lo"

func contains(s []string, str string) bool {
	return lo.Contains(s, str)
}
