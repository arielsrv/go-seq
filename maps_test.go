package seq_test

import (
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/sectrean/go-seq/internal/seqtest"
	"github.com/sectrean/go-seq/internal/testtypes"
	"github.com/stretchr/testify/assert"
)

func Test_AggregateGrouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		initFunc func(string) int
		f        func(int, int) int
		expected map[string]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{"a": 6},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{"a": 4, "b": 2},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.AggregateGrouped(tt.seq, tt.initFunc, tt.f)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_CountGrouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		expected map[string]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			expected: map[string]int{"a": 3},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			expected: map[string]int{"a": 2, "b": 1},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.CountGrouped(tt.seq)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_CountFuncGrouped(t *testing.T) {
	tests := []struct {
		name      string
		seq       iter.Seq2[string, int]
		predicate func(string, int) bool
		expected  map[string]int
	}{
		{
			name: "count even values",
			seq:  seq.Zip(seq.Yield("a", "b", "a", "a"), seq.Yield(1, 2, 3, 4)),
			predicate: func(k string, v int) bool {
				return v%2 == 0
			},
			expected: map[string]int{"a": 1, "b": 1},
		},
		{
			name: "count all values",
			seq:  seq.Zip(seq.Yield("a", "b", "a", "b"), seq.Yield(1, 2, 3, 4)),
			predicate: func(k string, v int) bool {
				return true
			},
			expected: map[string]int{"a": 2, "b": 2},
		},
		{
			name: "empty sequence",
			seq:  seq.Empty2[string, int](),
			predicate: func(k string, v int) bool {
				return false
			},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.CountFuncGrouped(tt.seq, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_Grouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		expected map[string][]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			expected: map[string][]int{"a": {1, 2, 3}},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			expected: map[string][]int{"a": {1, 3}, "b": {2}},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			expected: map[string][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.Grouped(tt.seq)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_Join(t *testing.T) {
	tests := []struct {
		name     string
		posts    iter.Seq[*testtypes.Post]
		users    map[int]*testtypes.User
		expected []*testtypes.UserPost
	}{
		{
			name: "matching users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 1, Title: "Post 1", Body: "Body 1"},
				&testtypes.Post{ID: 2, UserID: 2, Title: "Post 2", Body: "Body 2"},
			),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
				2: {ID: 2, Name: "User 2"},
			},
			expected: []*testtypes.UserPost{
				{UserName: "User 1", Body: "Body 1", Title: "Post 1"},
				{UserName: "User 2", Body: "Body 2", Title: "Post 2"},
			},
		},
		{
			name: "non-matching users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 3, Title: "Post 1", Body: "Body 1"},
			),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
				2: {ID: 2, Name: "User 2"},
			},
			expected: nil,
		},
		{
			name:  "empty posts",
			posts: seq.Empty[*testtypes.Post](),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
			},
			expected: nil,
		},
		{
			name: "empty users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 1, Title: "Post 1", Body: "Body 1"},
			),
			users:    map[int]*testtypes.User{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.Join(tt.posts, tt.users,
				func(p *testtypes.Post) int { return p.UserID },
				func(p *testtypes.Post, u *testtypes.User) *testtypes.UserPost {
					return &testtypes.UserPost{UserName: u.Name, Body: p.Body, Title: p.Title}
				},
			)
			seqtest.AssertEqual(t, tt.expected, result)
		})
	}
}

func Test_OuterJoin(t *testing.T) {
	tests := []struct {
		name     string
		posts    iter.Seq[*testtypes.Post]
		users    map[int]*testtypes.User
		expected []*testtypes.UserPost
	}{
		{
			name: "matching users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 1, Title: "Post 1", Body: "Body 1"},
				&testtypes.Post{ID: 2, UserID: 2, Title: "Post 2", Body: "Body 2"},
			),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
				2: {ID: 2, Name: "User 2"},
			},
			expected: []*testtypes.UserPost{
				{UserName: "User 1", Body: "Body 1", Title: "Post 1"},
				{UserName: "User 2", Body: "Body 2", Title: "Post 2"},
			},
		},
		{
			name: "non-matching users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 3, Title: "Post 1", Body: "Body 1"},
			),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
				2: {ID: 2, Name: "User 2"},
			},
			expected: []*testtypes.UserPost{
				{UserName: "[Unknown]", Body: "Body 1", Title: "Post 1"},
			},
		},
		{
			name:  "empty posts",
			posts: seq.Empty[*testtypes.Post](),
			users: map[int]*testtypes.User{
				1: {ID: 1, Name: "User 1"},
			},
			expected: nil,
		},
		{
			name: "empty users",
			posts: seq.Yield(
				&testtypes.Post{ID: 1, UserID: 1, Title: "Post 1", Body: "Body 1"},
			),
			users: map[int]*testtypes.User{},
			expected: []*testtypes.UserPost{
				{UserName: "[Unknown]", Body: "Body 1", Title: "Post 1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.OuterJoin(tt.posts, tt.users,
				func(p *testtypes.Post) int { return p.UserID },
				func(p *testtypes.Post, u *testtypes.User, found bool) *testtypes.UserPost {
					username := "[Unknown]"
					if found {
						username = u.Name
					}

					return &testtypes.UserPost{UserName: username, Body: p.Body, Title: p.Title}
				},
			)
			seqtest.AssertEqual(t, tt.expected, result)
		})
	}
}
