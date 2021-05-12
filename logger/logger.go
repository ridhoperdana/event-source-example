package logger

import "context"

type Logger interface {
	Process(ctx context.Context)
}

type Storage interface {
	Store(ctx context.Context, data []byte) error
}

type Service interface {
	Process(ctx context.Context, payload []byte) ([]byte, error)
}
