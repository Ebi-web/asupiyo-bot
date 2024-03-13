package event

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"log"
)

type ChickenEvent struct {
	bot *messaging_api.MessagingApiAPI
}

func NewChickenEvent(bot *messaging_api.MessagingApiAPI) *ChickenEvent {
	return &ChickenEvent{
		bot: bot,
	}
}

func (ce *ChickenEvent) IsApplicable(text string) bool {
	return text == "たまご" || text == "ひよこ" || text == "にわとり"
}

func (ce *ChickenEvent) Handle(tmc webhook.TextMessageContent, e webhook.MessageEvent) error {
	if _, err := ce.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: "ピヨピヨ",
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
