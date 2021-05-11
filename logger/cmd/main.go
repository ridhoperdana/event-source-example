package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/logger/internal"
)

func main() {
	redisClient := internal.NewRedisStorage(0, "", "localhost:6379")
	pubSub := internal.NewWatermillPubsub("event-source-example", "gateway-account", "account-balance",
		"gateway-account-sub", "pubsub-cred.json", redisClient)
	pubSub.Process(context.TODO())
}
