package robot

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"cyberstrike-ai/internal/config"

	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	dingutils "github.com/open-dingtalk/dingtalk-stream-sdk-go/utils"
	"go.uber.org/zap"
)

const (
	dingReconnectInitial = 5 * time.Second  // Initial reconnection interval
	dingReconnectMax     = 60 * time.Second // Maximum reconnection interval
)

// StartDing starts a DingTalk Stream persistent connection (no public endpoint required), invokes the handler for incoming messages, and replies through SessionWebhook.
// It reconnects automatically after disconnection (such as laptop sleep or network interruption) and exits when ctx is canceled, allowing restart after configuration changes.
func StartDing(ctx context.Context, robotsCfg config.RobotsConfig, h MessageHandler, logger *zap.Logger) {
	cfg := robotsCfg.Dingtalk
	if !cfg.Enabled || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return
	}
	go runDingLoop(ctx, cfg, robotsCfg.Session.StrictUserIdentityEnabled(), h, logger)
}

// runDingLoop maintains the DingTalk persistent connection and reconnects with backoff when disconnected while ctx remains active.
func runDingLoop(ctx context.Context, cfg config.RobotDingtalkConfig, strictUserIdentity bool, h MessageHandler, logger *zap.Logger) {
	backoff := dingReconnectInitial
	for {
		streamClient := client.NewStreamClient(
			client.WithAppCredential(client.NewAppCredentialConfig(cfg.ClientID, cfg.ClientSecret)),
			client.WithSubscription(dingutils.SubscriptionTypeKCallback, "/v1.0/im/bot/messages/get",
				chatbot.NewDefaultChatBotFrameHandler(func(ctx context.Context, msg *chatbot.BotCallbackDataModel) ([]byte, error) {
					go handleDingMessage(ctx, msg, cfg, strictUserIdentity, h, logger)
					return nil, nil
				}).OnEventReceived),
		)
		logger.Info("Connecting to DingTalk Stream...", zap.String("client_id", cfg.ClientID))
		err := streamClient.Start(ctx)
		if ctx.Err() != nil {
			logger.Info("DingTalk Stream closed for configuration restart")
			return
		}
		if err != nil {
			logger.Warn("DingTalk Stream persistent connection disconnected (for example, sleep or network outage); reconnecting automatically", zap.Error(err), zap.Duration("retry_after", backoff))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
			// Increase the next reconnection interval up to 60 seconds to avoid frequent retries
			if backoff < dingReconnectMax {
				backoff *= 2
				if backoff > dingReconnectMax {
					backoff = dingReconnectMax
				}
			}
		}
	}
}

func handleDingMessage(ctx context.Context, msg *chatbot.BotCallbackDataModel, cfg config.RobotDingtalkConfig, strictUserIdentity bool, h MessageHandler, logger *zap.Logger) {
	if msg == nil || msg.SessionWebhook == "" {
		return
	}
	content := ""
	if msg.Text.Content != "" {
		content = strings.TrimSpace(msg.Text.Content)
	}
	if content == "" && msg.Msgtype == "richText" {
		if cMap, ok := msg.Content.(map[string]interface{}); ok {
			if rich, ok := cMap["richText"].([]interface{}); ok {
				for _, c := range rich {
					if m, ok := c.(map[string]interface{}); ok {
						if txt, ok := m["text"].(string); ok {
							content = strings.TrimSpace(txt)
							break
						}
					}
				}
			}
		}
	}
	if content == "" {
		logger.Debug("Ignoring empty DingTalk message", zap.String("msgtype", msg.Msgtype))
		return
	}
	logger.Info("DingTalk message received", zap.String("sender", msg.SenderId), zap.String("content", content))
	tenantKey := strings.TrimSpace(cfg.ClientID)
	if tenantKey == "" {
		tenantKey = "default"
	}
	userID := strings.TrimSpace(msg.SenderId)
	if userID != "" {
		userID = "t:" + tenantKey + "|u:" + userID
	} else if cfg.AllowConversationIDFallback && !strictUserIdentity {
		conversationID := strings.TrimSpace(msg.ConversationId)
		if conversationID != "" {
			userID = "t:" + tenantKey + "|c:" + conversationID
		}
	}
	if userID == "" {
		logger.Warn("Ignoring DingTalk message without a usable user identifier")
		return
	}
	reply := h.HandleMessage("dingtalk", userID, content)
	// Use the markdown type to render headings, lists, code blocks, and similar formatting correctly
	title := reply
	if idx := strings.IndexAny(reply, "\n"); idx > 0 {
		title = strings.TrimSpace(reply[:idx])
	}
	if len(title) > 50 {
		title = title[:50] + "…"
	}
	if title == "" {
		title = "Reply"
	}
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  reply,
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, msg.SessionWebhook, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Warn("Failed to build DingTalk reply request", zap.Error(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Warn("DingTalk reply request failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Warn("DingTalk reply returned a non-200 status", zap.Int("status", resp.StatusCode))
		return
	}
	logger.Debug("DingTalk reply sent successfully", zap.String("content_preview", reply))
}
