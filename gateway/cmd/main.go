package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/ridhoperdana/event-source-example/gateway"

	"github.com/ridhoperdana/event-source-example/gateway/internal"
)

func main() {
	var (
		AccountID, EventType string
		Amount               int64
	)
	flag.StringVar(&AccountID, "account_id", "0", "ID of the account who wants to be processed")
	flag.StringVar(&EventType, "event_type", "INPUTTED_MONEY", "Event that want to be processed")
	flag.Int64Var(&Amount, "amount", 0, "amount of the processed event")

	flag.Parse()

	if Amount <= 0 {
		log.Fatal("please input amount > 0")
	}

	payload := gateway.ClientRequest{
		AccountID: AccountID,
		TypeRequest: gateway.TypeRequest{
			Type:   EventType,
			Amount: uint64(Amount),
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
