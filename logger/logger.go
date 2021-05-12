package logger

import (
	"context"

	"github.com/ridhoperdana/event-source-example/gateway"
)

type Logger interface {
	Process(ctx context.Context)
}

type Storage interface {
	Store(ctx context.Context, data Event) error
	StoreSnapshot(ctx context.Context, latestEvent Event) error
	Fetch(ctx context.Context, filter Filter) ([]Event, string, error)
	GetLatestEvent(ctx context.Context, id string) (Event, error)
	GetLatestSnapshot(ctx context.Context, id string) (AccountSnapshot, error)
}

type Service interface {
	Process(ctx context.Context, payload []byte) ([]byte, error)
}

type Filter struct {
	Num    int
	Cursor string
	ID     string
}

type Event struct {
	gateway.ClientRequest
	Version int `json:"version"`
}

type AccountSnapshot struct {
	gateway.Account
	Version int `json:"version"`
}
