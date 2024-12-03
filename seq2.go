package seq

import (
	"iter"
)

// Concat2 concatenates multiple key-value sequences into a single key-value sequence.
func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// ContainsKey determines whether a key is present in a key-value sequence.
func ContainsKey[K comparable, V any](seq iter.Seq2[K, V], key K) bool {
	for k, _ := range seq {
		if k == key {
			return true
		}
	}
	return false
}

// Empty2 returns an empty key-value sequence.
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}

// Keys returns the keys from a key-value sequence.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Select2 projects each key-value pair of a sequence into a new form.
func Select2[K, V, KOut, VOut any](seq iter.Seq2[K, V], f func(K, V) (KOut, VOut)) iter.Seq2[KOut, VOut] {
	return func(yield func(KOut, VOut) bool) {
		for k, v := range seq {
			kOut, vOut := f(k, v)
			if !yield(kOut, vOut) {
				return
			}
		}
	}
}

// SelectValues projects each key-value pair of a sequence into a new value.
func SelectValues[K, V, VOut any](seq iter.Seq2[K, V], f func(K, V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for k, v := range seq {
			out := f(k, v)
			if !yield(out) {
				return
			}
		}
	}
}

// Where2 filters key-value pairs based on a predicate.
func Where2[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// WithIndex returns a key-value sequence that incorporates the index of each value.
func WithIndex[V any](seq iter.Seq[V]) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Values returns the values from a key-value sequence.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Yield2 returns a sequence of indexes and values.
//
// This can be use to create a sequence from a slice or variadic arguments.
// The index of the value is used as the key.
//
// Examples:
//
//	// yield each element of the slice with its index
//	// yields (0, "a"), (1, "b"), (2, "c")
//	kvs := seq.Yield2([]string{"a", "b", "c"}...)
//
//	// yields (0, "a"), (1, "b"), (2, "c")
//	kvs := seq.Yield2("a", "b", "c")
func Yield2[V any](vals ...V) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range vals {
			if !yield(i, v) {
				return
			}
		}
	}
}

// Yield2Backward returns a sequence of indexes and values in reverse order.
//
// See Yield2 for more information.
func Yield2Backward[V any](vals ...V) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(vals) - 1; i >= 0; i-- {
			if !yield(i, vals[i]) {
				return
			}
		}
	}
}

// YieldMap returns a key-value sequence that yields the key-value pairs from a map.
//
// The iteration order is not specified and is not guaranteed to be the same from one call to the next.
func YieldMap[Map ~map[K]V, K comparable, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Zip returns a key-value sequence that combines values from two value sequences.
//
// The resulting sequence will be as long as the shortest value sequence.
// Any remaining values in the longer sequence are ignored.
//
// Example:
//
//	seqK := seq.Yield(1, 2, 3)
//	seqV := seq.Yield("a", "b", "c", "d")
//
//	// yields (1, "a"), (2, "b"), (3, "c")
//	kvs := seq.Zip(seqK, seqV)
func Zip[K any, V any](seqK iter.Seq[K], seqV iter.Seq[V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		nextV, stop := iter.Pull(seqV)
		defer stop()

		for k := range seqK {
			v, ok := nextV()
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
