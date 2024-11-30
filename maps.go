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

// Grouped collects a sequence of key-value pairs into a map with values grouped by key
// into slices.
func Grouped[K comparable, V any](seq iter.Seq2[K, V]) map[K][]V {
	groups := make(map[K][]V)
	for k, v := range seq {
		groups[k] = append(groups[k], v)
	}
	return groups
}
