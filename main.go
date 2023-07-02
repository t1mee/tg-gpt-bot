package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)


func goDotEnvVariable(key string) string {

  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func askChatGPT (question string) string {
	TOKEN := goDotEnvVariable("OPENAI_KEY")

	client := openai.NewClient(TOKEN)

    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    "user",
                    Content: question,
                },
            },
        },
    )

    if err != nil {
        panic(err)
    }

    return resp.Choices[0].Message.Content
}

func telegramBotProxy () {
	TOKEN := goDotEnvVariable("TG_KEY")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
    if err != nil {
        panic(err)
    }

    updateConfig := tgbotapi.NewUpdate(0)

    updateConfig.Timeout = 30

    updates := bot.GetUpdatesChan(updateConfig)

    for update := range updates {
        if update.Message == nil {
            continue
        }

		if update.Message.IsCommand() { // check if the message is a command
            switch update.Message.Command() {
            case "start":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я Snotra")
                bot.Send(msg)
				continue
            }
        }

		answer := askChatGPT(update.Message.Text)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

        msg.ReplyToMessageID = update.Message.MessageID

        if _, err := bot.Send(msg); err != nil {
            panic(err)
        }
    }

}

func main() {
	telegramBotProxy()
}