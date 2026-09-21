package robot

// MessageHandler is the message-processing interface called by Lark/DingTalk persistent connections and implemented by handler.RobotHandler.
type MessageHandler interface {
	HandleMessage(platform, userID, text string) string
}
