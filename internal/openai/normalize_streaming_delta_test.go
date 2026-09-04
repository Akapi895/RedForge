package openai

import "testing"

func TestNormalizeStreamingDelta_RepeatedCharBoundary(t *testing.T) {
	// A stream splits at a repeated-digit boundary: do not incorrectly merge the first character of "43" with the last character of "194".
	cur, d := normalizeStreamingDelta("https://x:194", "43")
	if want := "https://x:19443"; cur != want {
		t.Fatalf("next: want %q got %q", want, cur)
	}
	if d != "43" {
		t.Fatalf("delta: want %q got %q", "43", d)
	}
}

func TestNormalizeStreamingDelta_CumulativePrefix(t *testing.T) {
	cur, d := normalizeStreamingDelta("\u4eca\u5929", "\u4eca\u5929\u5929\u6c14")
	if cur != "\u4eca\u5929\u5929\u6c14" || d != "\u5929\u6c14" {
		t.Fatalf("got cur=%q d=%q", cur, d)
	}
}

func TestNormalizeStreamingDelta_FullRetransmit(t *testing.T) {
	cur, d := normalizeStreamingDelta("\u4eca\u5929", "\u4eca\u5929")
	if d != "" || cur != "\u4eca\u5929" {
		t.Fatalf("got cur=%q d=%q", cur, d)
	}
}

func TestNormalizeStreamingDelta_SingleRuneRepeated(t *testing.T) {
	cur, d := normalizeStreamingDelta("\u5440", "\u5440")
	if want := "\u5440\u5440"; cur != want {
		t.Fatalf("next: want %q got %q", want, cur)
	}
	if d != "\u5440" {
		t.Fatalf("delta: want %q got %q", "\u5440", d)
	}
	cur, d = normalizeStreamingDelta("4", "4")
	if want := "44"; cur != want {
		t.Fatalf("next: want %q got %q", want, cur)
	}
	if d != "4" {
		t.Fatalf("delta: want %q got %q", "4", d)
	}
}

func TestNormalizeStreamingDelta_CumulativeExtendsNumber(t *testing.T) {
	// After buffering "194", receive the accumulated string "19443" (note that "1943" is not a prefix of "19443", so HasPrefix must not be tested with an incorrectly written intermediate state).
	cur, d := normalizeStreamingDelta("194", "19443")
	if want := "19443"; cur != want {
		t.Fatalf("next: want %q got %q", want, cur)
	}
	if d != "43" {
		t.Fatalf("delta: want %q got %q", "43", d)
	}
}
