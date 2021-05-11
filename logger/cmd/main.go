package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/logger/internal"
)

func main() {
	redisClient := internal.NewRedisStorage(0, "", "localhost:6379")
	watermillMessenger := internal.NewWatermillPublisher("event-source-example", "account-balance", "pubsub-cred.json")
	pubSub := internal.NewWatermillPubsub("event-source-example", "gateway-account",
		"gateway-account-sub", "pubsub-cred.json", redisClient, watermillMessenger)
	pubSub.Process(context.TODO())
}
