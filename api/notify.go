package api

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func Notify(notification string) error {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	_, err := client.LPush(ctx, "TgNotifications", notification).Result()
	if err != nil {
		return err
	}
	return nil
}
