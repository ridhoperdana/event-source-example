package logger

import "context"

type Logger interface {
	Process(ctx context.Context)
}

type Storage interface {
	Store(ctx context.Context, data []byte) error
}
