package seq

import (
	"iter"
	"maps"
)

// func FromMap[Map ~map[K]V, K comparable, V any](m Map) iter.Seq2[K, V] {
// 	return maps.All(m)
// }

// func FromMapKeys[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[K] {
// 	return maps.Keys(m)
// }

// func FromMapValues[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[V] {
// 	return maps.Values(m)
// }

func CollectMap[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	return maps.Collect(seq)
}

func Grouped[K comparable, V any](seq iter.Seq2[K, V]) map[K][]V {
	groups := make(map[K][]V)
	for k, v := range seq {
		groups[k] = append(groups[k], v)
	}
	return groups
}
