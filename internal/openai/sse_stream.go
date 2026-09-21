package openai

// SSEAccumulatedKey is the server-authoritative full streaming snapshot field in SSE progress-event data.
// The frontend should prefer this field when updating its buffer to avoid duplicated characters from normalizing a delta twice.
const SSEAccumulatedKey = "accumulated"

// WithSSEAccumulated adds the current accumulated streaming text (authoritative snapshot) to progress data.
func WithSSEAccumulated(data map[string]interface{}, accumulated string) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{}, 1)
	}
	data[SSEAccumulatedKey] = accumulated
	return data
}

// NormalizeStreamingDelta normalizes content that may be an accumulated or retransmitted fragment into a pure delta.
// It matches the unexported normalizeStreamingDelta and is used by agent, multiagent, and similar packages to accumulate body text before sending SSE.
func NormalizeStreamingDelta(current, incoming string) (next, delta string) {
	return normalizeStreamingDelta(current, incoming)
}
