package project

import "strings"

// FactIndexSectionHeading is the readable heading prefix for the blackboard index (retained in the block for the Agent).
const FactIndexSectionHeading = "## Project Blackboard Index"

// FactIndexSectionStartMarker / EndMarker are HTML comment boundaries used for programmatic replacement and carry no instructional meaning for the model.
const (
	FactIndexSectionStartMarker = "<!-- fact-index-start -->"
	FactIndexSectionEndMarker   = "<!-- fact-index-end -->"
)

// ReplaceFactIndexSection replaces the existing project-blackboard index section in content with freshIndex.
// freshIndex must be the complete output of BuildFactIndexBlock. It returns (_, false) when either HTML boundary comment is missing.
func ReplaceFactIndexSection(content, freshIndex string) (string, bool) {
	freshIndex = strings.TrimSpace(freshIndex)
	if freshIndex == "" {
		return content, false
	}
	start, ok := factIndexSectionStart(content)
	if !ok {
		return content, false
	}
	end, ok := factIndexSectionEnd(content, start)
	if !ok || end <= start {
		return content, false
	}
	return content[:start] + freshIndex + content[end:], true
}

// wrapFactIndexBlock adds consistent opening and closing HTML comment boundaries around the BuildFactIndexBlock body.
func wrapFactIndexBlock(content string) string {
	content = strings.TrimSpace(content)
	return FactIndexSectionStartMarker + "\n" + content + "\n" + FactIndexSectionEndMarker + "\n"
}

func factIndexSectionStart(content string) (int, bool) {
	idx := strings.Index(content, FactIndexSectionStartMarker)
	if idx < 0 {
		return 0, false
	}
	return idx, true
}

func factIndexSectionEnd(content string, start int) (int, bool) {
	if start < 0 || start >= len(content) {
		return 0, false
	}
	tail := content[start:]
	idx := strings.LastIndex(tail, FactIndexSectionEndMarker)
	if idx < 0 {
		return 0, false
	}
	return start + idx + len(FactIndexSectionEndMarker), true
}
