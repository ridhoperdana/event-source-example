package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/balance/internal"
)

func main() {
	//pubSub := internal.NewWatermillPubsub("event-source-example", "gateway-account",
	//	"gateway-account-sub", "pubsub-cred.json", redisClient, watermillMessenger)
	m := internal.NewWatermillPublisher("event-source-example", "gateway-account", "pubsub-cred.json")
	p := internal.NewWatermillPubsub("event-source-example", "account-balance",
		"account-balance-sub", "pubsub-cred.json", m)
	p.Process(context.TODO())
}
