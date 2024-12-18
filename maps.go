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

// CountFuncGrouped counts the number of occurrences of each key in a sequence of key-value pairs
// based on a predicate.
func CountFuncGrouped[K comparable, V any](
	seq iter.Seq2[K, V],
	f func(K, V) bool,
) map[K]int {
	groups := make(map[K]int)

	for k, v := range seq {
		if f(k, v) {
			groups[k]++
		}
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

// Join joins a sequence with values from a map using a function to select the key and a function
// to project the values to a new value.
//
// The resulting sequence will only contain values where the key exists in the map.
//
// Example:
//
//	 var posts iter.Seq[*post.Post]
//	 var users map[user.UserID]*user.User
//
//	 postUsers := seq.Join(posts, users,
//		func(p *post.Post) user.UserID { return p.UserID },
//		func(p *post.Post, u *user.User) *post.PostUser { return post.NewPostUser(p, u) },
//	 )
func Join[V1 any, K comparable, Map ~map[K]V2, V2 any, VOut any](
	seq iter.Seq[V1],
	m Map,
	selectKey func(V1) K,
	f func(V1, V2) VOut,
) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v1 := range seq {
			k := selectKey(v1)
			v2, ok := m[k]

			if !ok {
				continue
			}

			out := f(v1, v2)
			if !yield(out) {
				return
			}
		}
	}
}

// OuterJoin joins a sequence with values from a map using a function to select the key and a function
// to project the values to a new value.
//
// The resulting sequence will contain all values from the sequence. A boolean flag is provided to
// indicate if the key was found in the map.
func OuterJoin[V1 any, K comparable, Map ~map[K]V2, V2 any, VOut any](
	seq iter.Seq[V1],
	m Map,
	selectKey func(V1) K,
	f func(V1, V2, bool) VOut,
) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v1 := range seq {
			k := selectKey(v1)
			v2, ok := m[k]
			out := f(v1, v2, ok)

			if !yield(out) {
				return
			}
		}
	}
}
