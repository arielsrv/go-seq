package seq

import (
	"iter"
	"maps"
)

// CollectMap collects a sequence of key-value pairs into a map.
// If there are duplicate keys, the last value for the key is kept.
func CollectMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	return maps.Collect(seq)
}

// AggregateGrouped aggregates values from a sequence of key-value pairs into a map with accumulated
// values grouped by key.
func AggregateGrouped[K comparable, V, A any](
	seq iter.Seq2[K, V],
	initFunc func(K) A,
	f func(A, V) A,
) map[K]A {
	groups := make(map[K]A)
	for k, v := range seq {
		acc, ok := groups[k]
		if !ok {
			acc = initFunc(k)
		}
		groups[k] = f(acc, v)
	}
	return groups
}

// CountGrouped counts the number of occurrences of each key in a sequence of key-value pairs.
func CountGrouped[K comparable, V any](seq iter.Seq2[K, V]) map[K]int {
	groups := make(map[K]int)
	for k := range seq {
		groups[k]++
	}
	return groups
}

// Grouped collects a sequence of key-value pairs into a map with values grouped by key
// into slices.
func Grouped[K comparable, V any](seq iter.Seq2[K, V]) map[K][]V {
	groups := make(map[K][]V)
	for k, v := range seq {
		groups[k] = append(groups[k], v)
	}
	return groups
}

// Join joins a map with a sequence of key-value pairs into a new sequence of key-value pairs.
func Join[Map ~map[K]V1, K comparable, V1, V2, VOut any](
	m Map,
	seq iter.Seq2[K, V2],
	f func(K, V1, V2) VOut,
) iter.Seq2[K, VOut] {
	return func(yield func(K, VOut) bool) {
		for k, v2 := range seq {
			v1, ok := m[k]
			if !ok {
				continue
			}
			out := f(k, v1, v2)
			if !yield(k, out) {
				return
			}
		}
	}
}
