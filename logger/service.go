package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ridhoperdana/event-source-example/gateway"
)

type loggerService struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return loggerService{
		storage: storage,
	}
}

func (l loggerService) Process(ctx context.Context, payload []byte) ([]byte, error) {
	clientRequest := gateway.ClientRequest{}
	if err := json.Unmarshal(payload, &clientRequest); err != nil {
		return nil, err
	}

	latestEvent, err := l.storage.GetLatestEvent(ctx, clientRequest.AccountID)
	if err != nil {
		return nil, err
	}

	event := Event{
		ClientRequest: clientRequest,
		Version:       latestEvent.Version + 1,
	}

	//TODO this logic should be in marketplace
	if clientRequest.TypeRequest.Type == gateway.EventRequestedCheckout {
		log.Println("Requested Checkout: ", string(payload))
		snapshot, err := l.storage.GetLatestSnapshot(ctx, event.AccountID)
		if err != nil {
			return nil, err
		}
		if snapshot.Balance < event.ClientRequest.TypeRequest.Amount {
			return nil, fmt.Errorf("cannot process checkout, balance minus")
		}
	}

	if err := l.storage.Store(ctx, event); err != nil {
		return nil, err
	}

	if err := l.storage.StoreSnapshot(ctx, event); err != nil {
		return nil, err
	}

	switch clientRequest.TypeRequest.Type {
	case gateway.EventStoreMoney:
		log.Println("Stored Money: ", string(payload))
	case gateway.EventInputMoney:
		log.Println("Inputted Money: ", string(payload))
		return payload, nil
	case gateway.EventRequestedCheckout:
		return payload, nil
	case gateway.EventProcessedCheckout:
		log.Println("Processed Checkout: ", string(payload))
	default:
		log.Println("event not supported: ", clientRequest.TypeRequest.Type)
	}

	return nil, nil
}
