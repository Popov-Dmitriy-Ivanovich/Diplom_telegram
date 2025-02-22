package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var notificationsChatIds map[int64]bool

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				if update.Message.Text != os.Getenv("SECRET_KEY") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный ключ доступа")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					continue
				}
				notificationsChatIds[update.Message.Chat.ID] = true
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно, вам будут приходить обновления")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		}
	}()
	for {
		res, err := rdb.BRPop(ctx, 0, "TgNotifications").Result()
		if err != nil {
			log.Println(err)
		}

		for id, _ := range notificationsChatIds {
			for _, notification := range res {
				msg := tgbotapi.NewMessage(id, notification)
				bot.Send(msg)
			}
		}
	}

}
