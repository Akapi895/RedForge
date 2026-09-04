package vision

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/llm"
	"cyberstrike-ai/internal/openai"

	einoopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

// Client invokes the standalone Vision ChatModel for a single Generate call.
type Client struct {
	cfg    config.VisionConfig
	mainOA config.OpenAIConfig
}

// NewClient constructs a vision client.
func NewClient(visionCfg config.VisionConfig, mainOpenAI config.OpenAIConfig) *Client {
	return &Client{cfg: visionCfg, mainOA: mainOpenAI}
}

// Analyze sends image bytes to the vision-language model and returns a textual description.
func (c *Client) Analyze(ctx context.Context, img ImagePayload, question string) (string, error) {
	if len(img.Bytes) == 0 {
		return "", fmt.Errorf("empty image payload")
	}
	mime := strings.TrimSpace(img.MIMEType)
	if mime == "" {
		mime = "image/jpeg"
	}
	oa := c.cfg.OpenAICfgEffective(c.mainOA)
	if strings.TrimSpace(oa.APIKey) == "" {
		return "", fmt.Errorf("vision API key is empty (set vision.api_key or openai.api_key)")
	}
	if strings.TrimSpace(oa.Model) == "" {
		return "", fmt.Errorf("vision model is empty")
	}

	timeout := time.Duration(c.cfg.TimeoutSecondsEffective()) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	httpClient := &http.Client{
		Timeout: timeout + 15*time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: timeout + 10*time.Second,
		},
	}

	b64 := base64.StdEncoding.EncodeToString(img.Bytes)
	detail := schema.ImageURLDetailLow
	switch c.cfg.DetailEffective() {
	case "high":
		detail = schema.ImageURLDetailHigh
	case "auto":
		detail = schema.ImageURLDetailAuto
	}

	prompt := buildVisionPrompt(question)
	if llm.IsClaudeProvider(oa.Provider) {
		nativeModel, err := llm.NewClaudeAgenticModel(
			ctx,
			oa,
			httpClient,
			oa.MaxCompletionTokensEffective(),
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("vision native Claude model: %w", err)
		}
		resp, err := nativeModel.Generate(ctx, []*schema.AgenticMessage{{
			Role: schema.AgenticRoleTypeUser,
			ContentBlocks: []*schema.ContentBlock{
				schema.NewContentBlock(&schema.UserInputText{Text: prompt}),
				schema.NewContentBlock(&schema.UserInputImage{
					Base64Data: b64,
					MIMEType:   mime,
					Detail:     detail,
				}),
			},
		}})
		if err != nil {
			return "", fmt.Errorf("vision native Claude generate: %w", err)
		}
		content, _ := llm.AgenticText(resp)
		if strings.TrimSpace(content) == "" {
			return "", fmt.Errorf("vision model returned empty content")
		}
		return strings.TrimSpace(content), nil
	}

	httpClient = openai.NewEinoHTTPClient(&oa, httpClient)
	maxCompletionTokens := oa.MaxCompletionTokensEffective()
	modelCfg := &einoopenai.ChatModelConfig{
		APIKey:              oa.APIKey,
		BaseURL:             strings.TrimSuffix(oa.BaseURL, "/"),
		Model:               oa.Model,
		HTTPClient:          httpClient,
		MaxCompletionTokens: &maxCompletionTokens,
	}
	chatModel, err := einoopenai.NewChatModel(ctx, modelCfg)
	if err != nil {
		return "", fmt.Errorf("vision chat model: %w", err)
	}
	userMsg := &schema.Message{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{Type: schema.ChatMessagePartTypeText, Text: prompt},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						Base64Data: &b64,
						MIMEType:   mime,
					},
					Detail: detail,
				},
			},
		},
	}

	resp, err := chatModel.Generate(ctx, []*schema.Message{userMsg})
	if err != nil {
		return "", fmt.Errorf("vision generate: %w", err)
	}
	if resp == nil || strings.TrimSpace(resp.Content) == "" {
		return "", fmt.Errorf("vision model returned empty content")
	}
	return strings.TrimSpace(resp.Content), nil
}

func buildVisionPrompt(question string) string {
	q := strings.TrimSpace(question)
	if q == "" {
		q = "Provide a general description of the image, focusing on authorized security-testing details (visible text, forms, buttons, CAPTCHAs, error messages, and technology-stack clues)."
	}
	extra := ""
	if looksLikeCaptchaQuestion(q) {
		extra = "\nIf this is a CAPTCHA, output only the character sequence you recognize, without spaces, punctuation, or explanation. If it is illegible, state clearly that it cannot be recognized."
	}
	return `You are an authorized security-testing assistant. Answer the user's question based on the image, describing only what you can verify from it without inventing details.
User question: ` + q + extra
}

func looksLikeCaptchaQuestion(q string) bool {
	s := strings.ToLower(q)
	for _, kw := range []string{"\u9a8c\u8bc1\u7801", "captcha", "verification code", "verify code", "vcode", "\u56fe\u5f62\u7801"} {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return strings.Contains(s, "\u53ea\u8f93\u51fa") && (strings.Contains(s, "\u5b57\u7b26") || strings.Contains(s, "character"))
}
