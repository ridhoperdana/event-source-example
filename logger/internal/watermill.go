package internal

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ridhoperdana/event-source-example/logger"
	"google.golang.org/api/option"
)

type BalanceWatermill struct {
	topic        string
	projectID    string
	credFilePath string
	publisher    message.Publisher
}

func NewBalanceWatermill(topic, projectID, credFilePath string) BalanceWatermill {
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
		log.Fatalf("error init google pubsub publisher balance service connection: %v", err)
	}
	return BalanceWatermill{
		topic:        topic,
		projectID:    projectID,
		credFilePath: credFilePath,
		publisher:    publisher,
	}
}

type watermillHandler struct {
	topic            string
	projectID        string
	subscription     string
	subscriber       *googlecloud.Subscriber
	balanceWatermill BalanceWatermill
	loggerService    logger.Service
}

func NewWatermillHandler(projectID, topic, subscription, credFilePath string,
	balanceService BalanceWatermill, loggerService logger.Service) logger.Logger {
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

	return watermillHandler{
		topic:            topic,
		projectID:        projectID,
		subscription:     subscription,
		subscriber:       subscriber,
		balanceWatermill: balanceService,
		loggerService:    loggerService,
	}
}

func (p watermillHandler) Process(ctx context.Context) {
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
		p.balanceWatermill.topic,
		p.balanceWatermill.publisher,
		func(msg *message.Message) ([]*message.Message, error) {
			updatedPayload, err := p.loggerService.Process(ctx, msg.Payload)
			if err != nil {
				return nil, err
			}

			if updatedPayload == nil || len(updatedPayload) == 0 {
				return nil, nil
			}

			newMessage := message.NewMessage(watermill.NewUUID(), msg.Payload)
			return []*message.Message{newMessage}, nil
		},
	)

	if err := router.Run(context.TODO()); err != nil {
		log.Fatalf("error running router: %v", err)
	}
}
