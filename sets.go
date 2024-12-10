package seq

import (
	"iter"
)

// Set is a set of values.
type Set[V comparable] map[V]struct{}

// Add adds a value to the set.
// Returns true if the value was added, false if it was already present.
func (s Set[V]) Add(v V) bool {
	if _, ok := s[v]; ok {
		return false
	}

	s[v] = struct{}{}

	return true
}

// Remove removes a value from the set.
// Returns true if the value was removed, false if it was not present.
func (s Set[V]) Remove(v V) bool {
	if _, ok := s[v]; !ok {
		return false
	}

	delete(s, v)

	return true
}

// Contains determines whether a value is present in the set.
func (s Set[V]) Contains(v V) bool {
	_, ok := s[v]
	return ok
}

// Values returns a sequence of values in the set.
func (s Set[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// CollectSet collects values from a sequence into a new set.
func CollectSet[V comparable](seq iter.Seq[V]) Set[V] {
	set := make(Set[V])
	for v := range seq {
		set[v] = struct{}{}
	}

	return set
}

// NewSet creates a new set from values.
func NewSet[V comparable](vals ...V) Set[V] {
	set := make(Set[V], len(vals))
	for _, v := range vals {
		set[v] = struct{}{}
	}

	return set
}

// Distinct returns distinct values from a sequence.
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
func Distinct[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		set := NewSet[V]()

		for v := range seq {
			if set.Add(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// DistinctKeys returns distinct key-value pairs from a sequence by comparing the keys.
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
func DistinctKeys[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		set := NewSet[K]()

		for k, v := range seq {
			if set.Add(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Except returns values from a sequence that are not present in a set.
func Except[V comparable](seq iter.Seq[V], vals Set[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		set := NewSet[V]()

		for v := range seq {
			if !vals.Contains(v) && set.Add(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ExceptKeys returns key-value pairs from a sequence when the key is not present in a set.
func ExceptKeys[K comparable, V any](seq iter.Seq2[K, V], keys Set[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		set := NewSet[K]()

		for k, v := range seq {
			if !keys.Contains(k) && set.Add(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Intersect returns values from a sequence that are present in a set.
func Intersect[V comparable](seq iter.Seq[V], vals Set[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		set := NewSet[V]()

		for v := range seq {
			if vals.Contains(v) && set.Add(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// IntersectKeys returns key-value pairs from a sequence when the key is present in a set.
func IntersectKeys[K comparable, V any](seq iter.Seq2[K, V], keys Set[K]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		set := NewSet[K]()

		for k, v := range seq {
			if keys.Contains(k) && set.Add(k) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Union returns the set union of multiple sequences.
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
func Union[V comparable](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		set := NewSet[V]()

		for _, seq := range seqs {
			for v := range seq {
				if set.Add(v) {
					if !yield(v) {
						return
					}
				}
			}
		}
	}
}

// UnionKeys returns the set union of multiple sequences of key-value pairs by comparing keys.
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
func UnionKeys[K comparable, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		set := NewSet[K]()

		for _, s := range seqs {
			for k, v := range s {
				if set.Add(k) {
					if !yield(k, v) {
						return
					}
				}
			}
		}
	}
}
