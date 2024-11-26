package seq

import (
	"iter"
)

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

func ContainsKey[K comparable, V any](seq iter.Seq2[K, V], key K) bool {
	for k, _ := range seq {
		if k == key {
			return true
		}
	}
	return false
}

func Count2[K, V any](seq iter.Seq2[K, V]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

func Count2Func[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) int {
	count := 0
	for k, v := range seq {
		if f(k, v) {
			count++
		}
	}
	return count
}

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

func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}

func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

func NewSeq2[V any](vals ...V) iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range vals {
			if !yield(i, v) {
				return
			}
		}
	}
}

func SelectValues[K, V, VOut any](seq iter.Seq2[K, V], f func(K, V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

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

func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}
