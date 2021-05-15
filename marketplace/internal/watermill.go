package internal

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ridhoperdana/event-source-example/gateway"

	"github.com/ridhoperdana/event-source-example/logger"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"google.golang.org/api/option"
)

type WatermillPubsub struct {
	topic        string
	projectID    string
	subscription string
	subscriber   *googlecloud.Subscriber
	publisher    gateway.Messenger
}

func NewWatermillPubsub(projectID, topic, subscription, credFilePath string, messenger gateway.Messenger) logger.Logger {
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			ProjectID: projectID,
			ClientOptions: []option.ClientOption{
				option.WithCredentialsFile(credFilePath),
			},
		},
		watermill.NewStdLogger(true, false),
	)
	if err != nil {
		log.Fatalf("error init google pubsub subscription connection: %v", err)
	}
	return WatermillPubsub{
		topic:        topic,
		projectID:    projectID,
		subscription: subscription,
		subscriber:   subscriber,
		publisher:    messenger,
	}
}

func (p WatermillPubsub) Process(ctx context.Context) {
	log.Println("Listening to event....")

	messages, err := p.subscriber.Subscribe(ctx, p.topic)
	if err != nil {
		log.Fatalf("error listening to subscription: %v", err)
	}

	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		ev := gateway.ClientRequest{}
		if err := json.Unmarshal(msg.Payload, &ev); err != nil {
			log.Println("error unmarshal event before send back to logger: ", err)
			continue
		}

		ev.TypeRequest.Type = gateway.EventProcessedCheckout
		payloadUpdated, err := json.Marshal(ev)
		if err != nil {
			log.Println("error marshal updated event before send back to logger: ", err)
			continue
		}
		if err := p.publisher.Publish(payloadUpdated); err != nil {
			log.Println("error sending back event to logger: ", err)
			continue
		}
		msg.Ack()
	}
}

type WatermillPublisher struct {
	topic     string
	projectID string
	publisher message.Publisher
}

func NewWatermillPublisher(projectID, topic, credFilePath string) gateway.Messenger {
	publisher, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID: projectID,
			ClientOptions: []option.ClientOption{
				option.WithCredentialsFile(credFilePath),
			},
		},
		watermill.NewStdLogger(true, false),
	)
	if err != nil {
		log.Fatalf("error init google pubsub subscription connection: %v", err)
	}
	return WatermillPublisher{
		topic:     topic,
		projectID: projectID,
		publisher: publisher,
	}
}

func (p WatermillPublisher) Publish(payload []byte) error {
	log.Println("sending back event: ", string(payload))
	return p.publisher.Publish(p.topic, message.NewMessage(watermill.NewUUID(), payload))
}
