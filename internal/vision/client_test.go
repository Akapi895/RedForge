package vision

import "testing"

func TestLooksLikeCaptchaQuestion(t *testing.T) {
	if !looksLikeCaptchaQuestion("Recognize the CAPTCHA and output characters only") {
		t.Fatal("expected captcha hint")
	}
	if looksLikeCaptchaQuestion("Describe the login-page layout") {
		t.Fatal("expected non-captcha")
	}
}
