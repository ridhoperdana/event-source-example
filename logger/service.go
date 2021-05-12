package logger

import (
	"context"
	"encoding/json"
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
	if err := l.storage.Store(ctx, payload); err != nil {
		return nil, err
	}
	ev := gateway.ClientRequest{}
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, err
	}

	switch ev.TypeRequest.Type {
	case gateway.EventStoreMoney:
		log.Println("Stored Money: ", string(payload))
	case gateway.EventInputMoney:
		log.Println("Inputted Money: ", string(payload))
		return payload, nil
	default:
		log.Println("event not supported: ", ev.TypeRequest.Type)
	}

	return nil, nil
}
