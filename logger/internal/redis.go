package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ridhoperdana/event-source-example/gateway"
	"github.com/ridhoperdana/event-source-example/logger"
)

type redisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(db int, password, address string) logger.Storage {
	return redisStorage{
		Client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}
}

func (r redisStorage) Store(ctx context.Context, data []byte) error {
	timeNow := time.Now()
	req := gateway.ClientRequest{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	key := fmt.Sprintf("%v:%v", req.AccountID, timeNow.UnixNano())
	return r.Client.SetNX(ctx, key, data, time.Duration(0)).Err()
}
