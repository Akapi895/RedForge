package robot

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"cyberstrike-ai/internal/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"go.uber.org/zap"
)

const (
	larkReconnectInitial = 5 * time.Second  // Initial reconnection interval
	larkReconnectMax     = 60 * time.Second // Maximum reconnection interval
)

type larkTextContent struct {
	Text string `json:"text"`
}

// StartLark starts a Lark persistent connection (no public endpoint required), invokes the handler for incoming messages, and sends replies.
// It reconnects automatically after disconnection (such as laptop sleep or network interruption) and exits when ctx is canceled, allowing restart after configuration changes.
func StartLark(ctx context.Context, robotsCfg config.RobotsConfig, h MessageHandler, logger *zap.Logger) {
	cfg := robotsCfg.Lark
	if !cfg.Enabled || cfg.AppID == "" || cfg.AppSecret == "" {
		return
	}
	go runLarkLoop(ctx, cfg, robotsCfg.Session.StrictUserIdentityEnabled(), h, logger)
}

// runLarkLoop maintains the Lark persistent connection and reconnects with backoff when disconnected while ctx remains active.
func runLarkLoop(ctx context.Context, cfg config.RobotLarkConfig, strictUserIdentity bool, h MessageHandler, logger *zap.Logger) {
	backoff := larkReconnectInitial
	for {
		larkClient := lark.NewClient(cfg.AppID, cfg.AppSecret)
		eventHandler := dispatcher.NewEventDispatcher("", "").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			go handleLarkMessage(ctx, event, cfg, strictUserIdentity, h, larkClient, logger)
			return nil
		})
		wsClient := larkws.NewClient(cfg.AppID, cfg.AppSecret,
			larkws.WithEventHandler(eventHandler),
			larkws.WithLogLevel(larkcore.LogLevelInfo),
		)
		logger.Info("Connecting to Lark persistent connection...", zap.String("app_id", cfg.AppID))
		err := wsClient.Start(ctx)
		if ctx.Err() != nil {
			logger.Info("Lark persistent connection closed for configuration restart")
			return
		}
		if err != nil {
			logger.Warn("Lark persistent connection disconnected (for example, sleep or network outage); reconnecting automatically", zap.Error(err), zap.Duration("retry_after", backoff))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
			if backoff < larkReconnectMax {
				backoff *= 2
				if backoff > larkReconnectMax {
					backoff = larkReconnectMax
				}
			}
		}
	}
}

func handleLarkMessage(ctx context.Context, event *larkim.P2MessageReceiveV1, cfg config.RobotLarkConfig, strictUserIdentity bool, h MessageHandler, client *lark.Client, logger *zap.Logger) {
	if event == nil || event.Event == nil || event.Event.Message == nil || event.Event.Sender == nil || event.Event.Sender.SenderId == nil {
		return
	}
	msg := event.Event.Message
	msgType := larkcore.StringValue(msg.MessageType)
	if msgType != larkim.MsgTypeText {
		logger.Debug("Lark currently processes text messages only", zap.String("msg_type", msgType))
		return
	}
	var textBody larkTextContent
	if err := json.Unmarshal([]byte(larkcore.StringValue(msg.Content)), &textBody); err != nil {
		logger.Warn("Failed to parse Lark message content", zap.Error(err))
		return
	}
	text := strings.TrimSpace(textBody.Text)
	if text == "" {
		return
	}
	userID := resolveLarkUserID(event, cfg.AllowChatIDFallback && !strictUserIdentity)
	if userID == "" {
		logger.Warn("Ignoring Lark message without a usable user identifier")
		return
	}
	messageID := larkcore.StringValue(msg.MessageId)
	reply := h.HandleMessage("lark", userID, text)
	contentBytes, _ := json.Marshal(larkTextContent{Text: reply})
	_, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(messageID).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Content(string(contentBytes)).
			Build()).
		Build())
	if err != nil {
		logger.Warn("Failed to send Lark reply", zap.String("message_id", messageID), zap.Error(err))
		return
	}
	logger.Debug("Lark reply sent", zap.String("message_id", messageID))
}

// resolveLarkUserID extracts the Lark conversation-isolation key:
// tenant_key plus a stable user identifier (user_id/open_id/union_id), with an optional configured chat_id fallback.
func resolveLarkUserID(event *larkim.P2MessageReceiveV1, allowChatIDFallback bool) string {
	if event == nil || event.Event == nil || event.Event.Sender == nil || event.Event.Sender.SenderId == nil {
		return ""
	}
	tenantKey := strings.TrimSpace(larkcore.StringValue(event.Event.Sender.TenantKey))
	if tenantKey == "" {
		tenantKey = "default"
	}
	prefix := "t:" + tenantKey + "|"
	if id := strings.TrimSpace(larkcore.StringValue(event.Event.Sender.SenderId.UserId)); id != "" {
		return prefix + "u:" + id
	}
	if id := strings.TrimSpace(larkcore.StringValue(event.Event.Sender.SenderId.OpenId)); id != "" {
		return prefix + "o:" + id
	}
	if id := strings.TrimSpace(larkcore.StringValue(event.Event.Sender.SenderId.UnionId)); id != "" {
		return prefix + "n:" + id
	}
	if allowChatIDFallback && event.Event.Message != nil {
		if id := strings.TrimSpace(larkcore.StringValue(event.Event.Message.ChatId)); id != "" {
			return prefix + "c:" + id
		}
	}
	return ""
}
