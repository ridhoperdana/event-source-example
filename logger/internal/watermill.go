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

type MarketplaceWatermill struct {
	handlerTopic      string
	handlerSubscriber *googlecloud.Subscriber
	targetTopic       string
	projectID         string
	credFilePath      string
	publisher         message.Publisher
}

func NewMarketplaceWatermill(targetTopic, projectID, credFilePath, handlerTopic string) MarketplaceWatermill {
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
		log.Fatalf("error init google pubsub publisher marketplace service connection: %v", err)
	}

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
		log.Fatalf("error init google pubsub subscription marketplace connection: %v", err)
	}
	return MarketplaceWatermill{
		handlerTopic:      handlerTopic,
		handlerSubscriber: subscriber,
		targetTopic:       targetTopic,
		projectID:         projectID,
		credFilePath:      credFilePath,
		publisher:         publisher,
	}
}

type BalanceWatermill struct {
	handlerTopic      string
	handlerSubscriber *googlecloud.Subscriber
	targetTopic       string
	projectID         string
	credFilePath      string
	publisher         message.Publisher
}

func NewBalanceWatermill(targetTopic, projectID, credFilePath, handlerTopic string) BalanceWatermill {
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
		log.Fatalf("error init google pubsub subscription balance connection: %v", err)
	}
	return BalanceWatermill{
		handlerTopic:      handlerTopic,
		handlerSubscriber: subscriber,
		targetTopic:       targetTopic,
		projectID:         projectID,
		credFilePath:      credFilePath,
		publisher:         publisher,
	}
}

type watermillHandler struct {
	projectID            string
	balanceWatermill     BalanceWatermill
	marketPlaceWatermill MarketplaceWatermill
	loggerService        logger.Service
}

func NewWatermillHandler(balanceService BalanceWatermill, marketPlaceWatermill MarketplaceWatermill,
	loggerService logger.Service) logger.Logger {

	return watermillHandler{
		balanceWatermill:     balanceService,
		loggerService:        loggerService,
		marketPlaceWatermill: marketPlaceWatermill,
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
		p.balanceWatermill.handlerTopic,
		p.balanceWatermill.handlerSubscriber,
		p.balanceWatermill.targetTopic,
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

	router.AddHandler(
		"handler account marketplace",
		p.marketPlaceWatermill.handlerTopic,
		p.marketPlaceWatermill.handlerSubscriber,
		p.marketPlaceWatermill.targetTopic,
		p.marketPlaceWatermill.publisher,
		func(msg *message.Message) ([]*message.Message, error) {
			log.Println("process marketplace")
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
