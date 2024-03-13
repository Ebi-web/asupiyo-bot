package event

import (
	"log"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type RajangEvent struct {
	bot *messaging_api.MessagingApiAPI
}

func NewRajangEvent(bot *messaging_api.MessagingApiAPI) *RajangEvent {
	return &RajangEvent{
		bot: bot,
	}
}

func (ce *RajangEvent) IsApplicable(text string) bool {
	return text == "ラージャン" || text == "らーじゃん"
}

func (ce *RajangEvent) Handle(tmc webhook.TextMessageContent, e webhook.MessageEvent) error {
	if _, err := ce.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.ImageMessage{
					PreviewImageUrl:    "https://www.monsterhunter.com/world-iceborne/ja/topics/e-jang/images/img_rajang01_l.png",
					OriginalContentUrl: "https://www.monsterhunter.com/world-iceborne/ja/topics/e-jang/images/img_rajang01_l.png",
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
