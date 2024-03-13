package event

import "github.com/line/line-bot-sdk-go/v8/linebot/webhook"

type TextMessageEvent interface {
	IsApplicable(text string) bool
	Handle(tmc webhook.TextMessageContent, e webhook.MessageEvent) error
}
