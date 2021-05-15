package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/marketplace/internal"
)

func main() {
	m := internal.NewWatermillPublisher("event-source-example", "gateway-account", "pubsub-cred.json")
	p := internal.NewWatermillPubsub("event-source-example", "account-marketplace",
		"account-marketplace-sub", "pubsub-cred.json", m)
	p.Process(context.TODO())
}
