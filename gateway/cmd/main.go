package main

import (
	"encoding/json"
	"log"

	"github.com/ridhoperdana/event-source-example/gateway"

	"github.com/ridhoperdana/event-source-example/gateway/internal"
)

func main() {
	payload := gateway.ClientRequest{
		AccountID: "3",
		TypeRequest: gateway.TypeRequest{
			Type:   gateway.EventInputMoney,
			Amount: 5000,
		},
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("error marshal payload: %v", err)
	}

	pub := internal.NewWatermillPubSub("event-source-example", "gateway-account",
		"pubsub-cred.json")

	if err := pub.Publish(payloadByte); err != nil {
		log.Fatalf("error publish message: %v", err)
	}
}
