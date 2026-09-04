package openai

// eino_sse_sanitizer.go addresses ErrTooManyEmptyStreamMessages when Eino uses the meguminnnnnnnnn/go-openai SDK
// and relay heartbeat/SSE control lines accumulate beyond 300
// (error message: "stream has sent too many empty messages").
//
// Trigger path:
//   einoopenai.NewChatModel
//     → eino-ext/libs/acl/openai → meguminnnnnnnnn/go-openai
//     → streamReader.processLines() counts every non-"data:" line and raises an error above 300.
//
// Common non-data lines from relays (valid SSE, but rejected by the SDK):
//   ":" / ": keepalive" / ": ping" / "event: ping" / "retry: 3000"
//   plus many interleaved heartbeats during reasoning-model prefill.
//
// Fallback strategy: wrap the response Body with a reader at the HTTP transport layer and pass through only lines
// beginning with "data:", discarding heartbeat, comment, and event-type lines in place. The downstream SDK never sees
// non-data lines, so its counter remains at zero and the error cannot occur.
//
// This layer is completely transparent to callers:
//   - It intervenes only when Content-Type is text/event-stream; ordinary JSON responses pass through unchanged.
//   - data payloads (including [DONE] and {"error":...}) are byte-for-byte unchanged.
//   - Genuine upstream termination (EOF, connection reset, or context cancellation) passes through unchanged.

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
)

const (
	// einoSSEReaderBufSize gives bufio a larger initial buffer so large single-line JSON chunks
	// (including tool-call arguments/reasoning_content) do not repeatedly trigger buffer growth.
	einoSSEReaderBufSize = 64 * 1024
)

// einoSSESanitizingRoundTripper wraps the downstream RoundTripper and sanitizes SSE responses line by line.
type einoSSESanitizingRoundTripper struct {
	base http.RoundTripper
}

func (rt *einoSSESanitizingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.base.RoundTrip(req)
	if err != nil || resp == nil {
		return resp, err
	}
	if !isSSEResponse(resp) {
		return resp, nil
	}
	resp.Body = newEinoSSESanitizingBody(resp.Body)
	return resp, nil
}

// isSSEResponse sanitizes only 200 + text/event-stream responses;
// error responses (4xx/5xx, usually application/json) remain unchanged so the SDK follows its original error path.
func isSSEResponse(resp *http.Response) bool {
	if resp.StatusCode != http.StatusOK {
		return false
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		return false
	}
	ct = strings.ToLower(strings.TrimSpace(ct))
	// Support "text/event-stream", "text/event-stream; charset=utf-8", and similar values.
	return strings.HasPrefix(ct, "text/event-stream")
}

// einoSSESanitizingBody is the wrapped response body: it passes through data lines and discards all others.
type einoSSESanitizingBody struct {
	upstream io.ReadCloser
	reader   *bufio.Reader
	pending  []byte // Sanitized bytes waiting to be returned downstream (always a complete data line ending in \n)
	err      error  // Terminal upstream error (io.EOF or network error)
}

func newEinoSSESanitizingBody(body io.ReadCloser) *einoSSESanitizingBody {
	return &einoSSESanitizingBody{
		upstream: body,
		reader:   bufio.NewReaderSize(body, einoSSEReaderBufSize),
	}
}

func (b *einoSSESanitizingBody) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(b.pending) > 0 {
		n := copy(p, b.pending)
		b.pending = b.pending[n:]
		return n, nil
	}

	// Read upstream until one data line is accumulated or a terminal state is reached.
	// One loop may discard any number of heartbeat lines, but exits after passing through at most one data line,
	// preventing a single Read from blocking too long or making the pending buffer too large.
	for b.err == nil {
		line, err := b.reader.ReadBytes('\n')
		if len(line) > 0 {
			if isPassThroughSSELine(line) {
				if line[len(line)-1] != '\n' {
					line = append(line, '\n')
				}
				b.pending = line
				if err != nil {
					b.err = err
				}
				break
			}
			// Discard every non-data line (blank lines, ":" comments, event:, retry:, id:, or any bare text)
			// without exposing it downstream, then continue reading the next line.
		}
		if err != nil {
			b.err = err
			break
		}
	}

	if len(b.pending) > 0 {
		n := copy(p, b.pending)
		b.pending = b.pending[n:]
		return n, nil
	}
	return 0, b.err
}

func (b *einoSSESanitizingBody) Close() error {
	return b.upstream.Close()
}

// isPassThroughSSELine determines whether a line should pass through unchanged to the downstream SDK.
// Only lines beginning with "data:" (case-insensitive, with arbitrary leading whitespace) are retained.
// Do not use TrimSpace to remove the trailing newline before testing, or "  data: x" may be misclassified;
// trim only leading whitespace, matching the SDK semantics of TrimSpace followed by the ^data:\s* regex.
func isPassThroughSSELine(line []byte) bool {
	trimmed := bytes.TrimLeft(line, " \t")
	if len(trimmed) < 5 {
		return false
	}
	// Compare the first five bytes with "data:" case-insensitively. The SSE specification requires lowercase field names,
	// but lenient matching supports non-standard implementations used by some relays.
	return (trimmed[0] == 'd' || trimmed[0] == 'D') &&
		(trimmed[1] == 'a' || trimmed[1] == 'A') &&
		(trimmed[2] == 't' || trimmed[2] == 'T') &&
		(trimmed[3] == 'a' || trimmed[3] == 'A') &&
		trimmed[4] == ':'
}
