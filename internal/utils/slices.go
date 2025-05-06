package utils

import (
	"cmp"
	"iter"
	"maps"
	"slices"
	"strings"
)

// SortSliceByOther sorts `slice` by the positioning of equal items that exist in the `other` slice.
//
// You can also provide `fallbackCompare` to be used if both items don't exist in the `other` slice. If only a single
// item is missing from the `other` slice, it will do a simple compare (a > b / a < b / a == b).
func SortSliceByOther[Slice ~[]E, E cmp.Ordered](slice, other Slice, fallbackCompare func(a, b E) int) Slice {
	positionsMap := maps.Collect(AllCrossed(other))
	return slices.SortedFunc(slices.Values(slice), CompareByMap(positionsMap, fallbackCompare))
}

func SortStringSliceByOther(slice, other []string) []string {
	return SortSliceByOther(slice, other, strings.Compare)
}

func Map[S, R any](items []S, fn func(S) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

func AllCrossed[Slice ~[]E, E comparable](s Slice) iter.Seq2[E, int] {
	return func(yield func(E, int) bool) {
		for i, v := range s {
			if !yield(v, i) {
				return
			}
		}
	}
}

func SimpleCompare[E cmp.Ordered](a, b E) int {
	switch true {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func CompareByMap[E cmp.Ordered](m map[E]int, fallbackCompare func(a, b E) int) func(a, b E) int {
	return func(a, b E) int {
		aPos, okA := m[a]
		bPos, okB := m[b]
		if !okA && !okB {
			if fallbackCompare == nil {
				return SimpleCompare(a, b)
			}
			return fallbackCompare(a, b)
		}
		return aPos - bPos
	}
}
