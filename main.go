package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"asupiyo-bot/event"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	bot, err := setupBot()
	if err != nil {
		log.Fatalf("Error setting up bot: %v", err)
	}

	textEvents := setupTextEvents(bot)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		handleCallback(w, req, bot, textEvents)
	})

	startServer()
}

func loadEnv() error {
	env := os.Getenv("ENV")
	if env == "" {
		fileName := ".env.local"
		if err := godotenv.Load(fileName); err != nil {
			return fmt.Errorf("error loading %s file: %w", fileName, err)
		}
	}
	return nil
}

func setupBot() (*messaging_api.MessagingApiAPI, error) {
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}
	return bot, nil
}

func setupTextEvents(bot *messaging_api.MessagingApiAPI) []event.TextMessageEvent {
	chickenEvent := event.NewChickenEvent(bot)
	defaultEvent := event.NewDefaultEvent(bot)
	return []event.TextMessageEvent{chickenEvent, defaultEvent}
}

func handleCallback(w http.ResponseWriter, req *http.Request, bot *messaging_api.MessagingApiAPI, textEvents []event.TextMessageEvent) {
	log.Println("/callback called...")

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	cb, err := webhook.ParseRequest(channelSecret, req)
	if err != nil {
		log.Printf("Cannot parse request: %+v\n", err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, webhook.ErrInvalidSignature) {
			statusCode = http.StatusBadRequest
		}
		w.WriteHeader(statusCode)
		return
	}

	processEvents(cb.Events, bot, textEvents)
}

func processEvents(events []webhook.EventInterface, bot *messaging_api.MessagingApiAPI, textEvents []event.TextMessageEvent) {
	log.Println("Handling events...")
	for _, callbackRequest := range events {
		handleEvent(callbackRequest, bot, textEvents)
	}
}

func handleEvent(callbackRequest webhook.EventInterface, bot *messaging_api.MessagingApiAPI, textEvents []event.TextMessageEvent) {
	switch e := callbackRequest.(type) {
	case webhook.MessageEvent:
		handleMessageEvent(e, bot, textEvents)
	default:
		log.Printf("Unsupported message: %T\n", callbackRequest)
	}
}

func handleMessageEvent(e webhook.MessageEvent, bot *messaging_api.MessagingApiAPI, textEvents []event.TextMessageEvent) {
	switch message := e.Message.(type) {
	case webhook.TextMessageContent:
		for _, textEvent := range textEvents {
			if textEvent.IsApplicable(message.Text) {
				textEvent.Handle(message, e) // Improved error handling is suggested
				break
			}
		}
	case webhook.StickerMessageContent:
		sendStickerReply(message, e.ReplyToken, bot)
	default:
		log.Printf("Unsupported message content: %T\n", e.Message)
	}
}

func sendStickerReply(message webhook.StickerMessageContent, replyToken string, bot *messaging_api.MessagingApiAPI) {
	replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
	if _, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: replyMessage,
			},
		},
	}); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent sticker reply.")
	}
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	log.Printf("Server starting on http://localhost:%s/\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
