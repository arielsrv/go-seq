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

// ContainsKey returns true if the given key is present in the key-value sequence.
func ContainsKey[K comparable, V any](seq iter.Seq2[K, V], key K) bool {
	for k, _ := range seq {
		if k == key {
			return true
		}
	}
	return false
}

// DistinctKeys returns distinct key-value pairs from a sequence by comparing keys.
func DistinctKeys[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seen := make(map[K]struct{})
		for k, v := range seq {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(k, v) {
					return
				}
			}
		}
	}
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

// NewSeq2 returns a key-value sequence that yields the given values and their indices.
func NewSeq2[V any](vals ...V) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range vals {
			if !yield(i, v) {
				return
			}
		}
	}
}

// SelectValues projects each key-value pair of a sequence into a new value.
func SelectValues[K, V, VOut any](seq iter.Seq2[K, V], f func(K, V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
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

// Zip returns a key-value sequence that combines values from two input sequences.
//
// The resulting sequence will be as long as the shortest input sequence.
// Any remaining values in the longer sequence are discarded.
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
