package pubsub

import (
	"strings"
)

// MatchTopic checks if a target topic matches a subscription pattern.
// Supports MQTT-style wildcards:
// '+' (or '*') matches a single level
// '#' matches zero or more levels at the end
func MatchTopic(pattern, target string) bool {
	if pattern == target {
		return true
	}

	pParts := strings.Split(pattern, ".")
	tParts := strings.Split(target, ".")

	return matchParts(pParts, tParts)
}

func matchParts(pattern, target []string) bool {
	if len(pattern) == 0 {
		return len(target) == 0
	}

	p := pattern[0]

	if p == "#" {
		return true
	}

	if len(target) == 0 {
		return false
	}

	if p == "+" || p == "*" || p == target[0] {
		return matchParts(pattern[1:], target[1:])
	}

	return false
}

// GetMatchingTopics returns all topics in the registry that match the pattern
func (r *TopicRegistry) GetMatchingTopics(pattern string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []string
	for topicName := range r.topics {
		if MatchTopic(pattern, topicName) {
			matches = append(matches, topicName)
		}
	}
	return matches
}
