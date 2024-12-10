package seqtest

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertEqual asserts that the expected values are equal to the actual values from a sequence.
func AssertEqual[V any](t *testing.T, expected []V, seq iter.Seq[V]) {
	t.Helper()

	actual := slices.Collect(seq)
	assert.Equal(t, expected, actual)
}

type KeyValuePair[K, V any] struct {
	Key   K
	Value V
}

// AssertEqual2 asserts that the expected key-value pairs are equal to the actual key-value pairs from a sequence.
func AssertEqual2[K, V any](t *testing.T, expected []KeyValuePair[K, V], seq iter.Seq2[K, V]) {
	t.Helper()

	//nolint:prealloc // We don't know the length of the sequence.
	var actual []KeyValuePair[K, V]
	for k, v := range seq {
		actual = append(actual, KeyValuePair[K, V]{k, v})
	}

	assert.Equal(t, expected, actual)
}
