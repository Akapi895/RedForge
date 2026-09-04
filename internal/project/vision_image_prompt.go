package project

import "strings"

// VisionImageSectionMarker is the image-analysis section heading used by AppendVisionImageAnalysisIfReady.
const VisionImageSectionMarker = "## Image Analysis"

// VisionImageAnalysisSection is the image-analysis guidance shared by single- and multi-agent modes (analyze_image; only the textual summary remains in context).
func VisionImageAnalysisSection() string {
	var b strings.Builder
	b.WriteString(VisionImageSectionMarker)
	b.WriteString("\n\n")
	b.WriteString("- When encountering an image file (screenshot, CAPTCHA, login page, or report illustration), pass its server-side file path to analyze_image when that tool is available.\n")
	b.WriteString("- Do not use read_file on a binary image and expect its contents to be understood; a user-message attachment in the form '📎 xxx.png: /path' provides a path that can be passed to analyze_image.\n")
	b.WriteString("- For CAPTCHAs saved locally from a page or endpoint (for example, captcha.png), use analyze_image and state in question: 'Output only the CAPTCHA characters.' If recognition fails, refresh and save the CAPTCHA again before retrying. Do not expect one image-analysis attempt to solve complex slider or behavioral CAPTCHAs.\n")
	b.WriteString("- When delegating to a sub-agent, if the task includes CAPTCHA or screenshot interpretation, include the image path and expected output format in the task description.\n")
	return b.String()
}

// AppendVisionImageAnalysisIfReady appends image-analysis guidance only when vision.enabled is true and a model is configured.
func AppendVisionImageAnalysisIfReady(base string, visionReady bool) string {
	if !visionReady {
		return base
	}
	return AppendSystemPromptBlock(base, VisionImageAnalysisSection())
}
