package project

import "cyberstrike-ai/internal/projectprompt"

// FactRecordingIncrementalRhythmMarkdown is documented in projectprompt.
func FactRecordingIncrementalRhythmMarkdown(coordinator, subAgent bool) string {
	return projectprompt.FactRecordingIncrementalRhythmMarkdown(coordinator, subAgent)
}

// FactRecordingBlackboardSection is documented in projectprompt.
func FactRecordingBlackboardSection(coordinatorDelegate bool) string {
	return projectprompt.FactRecordingBlackboardSection(coordinatorDelegate)
}

// FactRecordingSubAgentSection is documented in projectprompt.
func FactRecordingSubAgentSection() string {
	return projectprompt.FactRecordingSubAgentSection()
}

// FactRecordingBlackboardSectionMarkdown is documented in projectprompt.
func FactRecordingBlackboardSectionMarkdown(coordinatorDelegate bool) string {
	return projectprompt.FactRecordingBlackboardSectionMarkdown(coordinatorDelegate)
}
