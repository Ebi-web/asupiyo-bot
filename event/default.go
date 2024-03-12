package event

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"log"
)

type DefaultEvent struct {
	bot *messaging_api.MessagingApiAPI
}

func NewDefaultEvent(bot *messaging_api.MessagingApiAPI) *DefaultEvent {
	return &DefaultEvent{
		bot: bot,
	}
}

func (ce *DefaultEvent) IsApplicable(text string) bool {
	return true
}

func (ce *DefaultEvent) Handle(tmc webhook.TextMessageContent, e webhook.MessageEvent) error {
	if _, err := ce.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: tmc.Text,
				},
			},
		},
	); err != nil {
		log.Print(err)
		return nil
	} else {
		log.Println("Sent text reply.")
		return err
	}
}
