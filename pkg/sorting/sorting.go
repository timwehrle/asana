package sorting

import (
	"sort"

	"cmp"
)

func By[T any](slice []T, less func(a T, b T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
}

func ByField[T any, F cmp.Ordered](slice []T, field func(T) F, descending bool) {
	sort.Slice(slice, func(i, j int) bool {
		if descending {
			return field(slice[i]) > field(slice[j])
		}
		return field(slice[i]) < field(slice[j])
	})
}

func Sort[T any](slice []T, by func(a, b T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return by(slice[i], slice[j])
	})
}
