package internal

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ridhoperdana/event-source-example/gateway"
	"github.com/ridhoperdana/event-source-example/logger"
	"google.golang.org/api/option"
)

type WatermillPubsub struct {
	topic        string
	topicBalance string
	projectID    string
	subscription string
	subscriber   *googlecloud.Subscriber
	storage      logger.Storage
	messenger    message.Publisher
}

func NewWatermillPubsub(projectID, topic, topicBalance, subscription, credFilePath string, storage logger.Storage) logger.Logger {
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
		log.Fatalf("error init google pubsub publisher connection: %v", err)
	}
	return WatermillPubsub{
		topic:        topic,
		projectID:    projectID,
		subscription: subscription,
		subscriber:   subscriber,
		storage:      storage,
		messenger:    publisher,
		topicBalance: topicBalance,
	}
}

func (p WatermillPubsub) Process(ctx context.Context) {
	watermillLogger := watermill.NewStdLogger(true, false)
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		log.Fatalf("error creating router: %v", err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.Recoverer)

	router.AddHandler(
		"handler account balance",
		p.topic,
		p.subscriber,
		p.topicBalance,
		p.messenger,
		func(msg *message.Message) ([]*message.Message, error) {
			log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
			if err := p.storage.Store(ctx, msg.Payload); err != nil {
				return nil, err
			}
			ev := gateway.ClientRequest{}
			if err := json.Unmarshal(msg.Payload, &ev); err != nil {
				return nil, err
			}

			switch ev.TypeRequest.Type {
			case gateway.EventStoreMoney:
				log.Println("Stored Money: ", string(msg.Payload))
				return nil, nil
			case gateway.EventInputMoney:
			default:
				log.Println("event not supported: ", ev.TypeRequest.Type)
				return nil, nil
			}

			newMessage := message.NewMessage(watermill.NewUUID(), msg.Payload)
			return []*message.Message{newMessage}, nil
		},
	)

	if err := router.Run(ctx); err != nil {
		log.Fatalf("error running router: %v", err)
	}
}
