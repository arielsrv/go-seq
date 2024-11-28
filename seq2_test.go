package seq_test

import (
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/sectrean/go-seq/internal/seqtest"
)

func Test_Zip(t *testing.T) {
	seqK := seq.Yield(1, 2, 3)
	seqV := seq.Yield("a", "b", "c", "d")

	result := seq.Zip(seqK, seqV)

	expected := []seqtest.KeyValuePair[int, string]{
		{Key: 1, Value: "a"},
		{Key: 2, Value: "b"},
		{Key: 3, Value: "c"},
	}
	seqtest.AssertEqual2(t, expected, result)
}
