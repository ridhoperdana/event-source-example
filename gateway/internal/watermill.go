package internal

import (
	"log"

	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ridhoperdana/event-source-example/gateway"
	"google.golang.org/api/option"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type WatermillPubsub struct {
	topic     string
	projectID string
	publisher message.Publisher
}

func NewWatermillPubSub(projectID, topic, credFilePath string) gateway.Messenger {
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
	return WatermillPubsub{
		topic:     topic,
		projectID: projectID,
		publisher: publisher,
	}
}

func (p WatermillPubsub) Publish(payload []byte) error {
	log.Println("sending event: ", string(payload))
	return p.publisher.Publish(p.topic, message.NewMessage(watermill.NewUUID(), payload))
}
