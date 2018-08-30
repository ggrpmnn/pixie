package pixie

// BotError types enable different messages for the log and the bot
type BotError struct {
	err    string
	botMsg string
}

// Error returns the technical error message, for logging purposes
func (b BotError) Error() string {
	return b.err
}

// BotMessage returns the friendly error message, for response purposes
func (b BotError) BotMessage() string {
	return b.botMsg
}
