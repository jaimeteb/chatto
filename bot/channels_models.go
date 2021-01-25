package bot

import "github.com/slack-go/slack"

// TelegramMessageIn models a telegram incoming message
type TelegramMessageIn struct {
	UpdateID int                    `json:"update_id"`
	Message  TelegramMessageInInner `json:"message"`
}

// TelegramMessageInInner models a telegram incoming message inner struct
type TelegramMessageInInner struct {
	MessageID int                        `json:"message_id"`
	From      TelegramMessageInInnerFrom `json:"from"`
	Date      int                        `json:"date"`
	Text      string                     `json:"text"`
}

// TelegramMessageInInnerFrom models a telegram incoming message inner struct
type TelegramMessageInInnerFrom struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// TwilioMessageIn models an incoming Twilio message
type TwilioMessageIn struct {
	From             string `form:"From"`
	Body             string `form:"Body"`
	To               string `form:"To"`
	MediaURL         string `form:"MediaUrl"`
	MediaContentType string `form:"MediaContentType"`
	MessageSid       string `form:"MessageSid"`
	SmsStatus        string `form:"SmsStatus"`
	AccountSid       string `form:"AccountSid"`
	Sid              string `form:"Sid"`
	SmsSid           string `form:"SmsSid"`
	SmsMessageSid    string `form:"SmsMessageSid"`
	NumMedia         int    `form:"NumMedia"`
	NumSegments      int    `form:"NumSegments"`
	APIVersion       string `form:"ApiVersion"`
}

// SlackMessage models a Slack message and/or Slack endpoint challenge
type SlackMessage struct {
	Challenge string    `json:"challenge"`
	Type      string    `json:"type"`
	Event     slack.Msg `json:"event"`
}
