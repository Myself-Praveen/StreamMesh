package pubsub

import (
	"testing"
)

func TestMatchTopic(t *testing.T) {
	tests := []struct {
		pattern  string
		target   string
		expected bool
	}{
		{"a.b.c", "a.b.c", true},
		{"a.b.c", "a.b.d", false},
		{"a.*.c", "a.b.c", true},
		{"a.*.c", "a.d.c", true},
		{"a.*.c", "a.b.d", false},
		{"a.#", "a.b.c", true},
		{"a.#", "a", true},
		{"#", "anything", true},
		{"a.b.*", "a.b", false},
		{"a.b.*", "a.b.c", true},
		{"a.b.*", "a.b.c.d", false},
	}

	for _, test := range tests {
		result := MatchTopic(test.pattern, test.target)
		if result != test.expected {
			t.Errorf("MatchTopic(%q, %q) = %v; expected %v", test.pattern, test.target, result, test.expected)
		}
	}
}
