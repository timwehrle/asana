package sorting

import (
	"sort"
)

func Sort[T any](slice []T, by func(a, b T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return by(slice[i], slice[j])
	})
}
