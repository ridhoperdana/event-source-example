package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func (r redisStorage) Store(ctx context.Context, data logger.Event) error {
	eventByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Client.LPush(ctx, data.AccountID, eventByte).Err()
}

func (r redisStorage) StoreSnapshot(ctx context.Context, latestEvent logger.Event) error {
	latestSnapshot, err := r.GetLatestSnapshot(ctx, latestEvent.AccountID)
	if err != nil {
		return err
	}

	if latestSnapshot.Version != latestEvent.Version && latestEvent.Version-latestSnapshot.Version != 6 {
		log.Println("snapshot not created, version not 6")
		return nil
	}

	previousEvent, _, err := r.Fetch(ctx, logger.Filter{
		ID: latestEvent.AccountID,
	})
	if err != nil {
		return err
	}

	if len(previousEvent) < 6 {
		log.Println("snapshot not created, version less than 6")
		return nil
	}

	snapshot := logger.AccountSnapshot{
		Account: gateway.Account{
			ID: latestEvent.AccountID,
		},
		Version: latestEvent.Version,
	}

	for _, event := range previousEvent {
		if event.TypeRequest.Type == gateway.EventStoreMoney {
			snapshot.Balance += event.TypeRequest.Amount
		}
	}

	snapshotByte, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	return r.Client.LPush(ctx, fmt.Sprintf("%v-snapshot", snapshot.ID), snapshotByte).Err()
}

func (r redisStorage) Fetch(ctx context.Context, filter logger.Filter) ([]logger.Event, string, error) {
	res, err := r.Client.LRange(ctx, filter.ID, 0, 9).Result()
	if err != nil {
		return nil, "", err
	}
	result := make([]logger.Event, len(res))
	for index, item := range res {
		var req logger.Event
		if err := json.Unmarshal([]byte(item), &req); err != nil {
			return nil, "", err
		}
		result[index] = req
	}

	return result, "", nil
}

func (r redisStorage) GetLatestEvent(ctx context.Context, id string) (logger.Event, error) {
	res, err := r.Client.LRange(ctx, id, 0, 1).Result()
	if err != nil {
		return logger.Event{}, err
	}
	if len(res) == 0 {
		return logger.Event{}, nil
	}
	var result logger.Event
	if err := json.Unmarshal([]byte(res[0]), &result); err != nil {
		return logger.Event{}, err
	}
	return result, nil
}

func (r redisStorage) GetLatestSnapshot(ctx context.Context, id string) (logger.AccountSnapshot, error) {
	id = fmt.Sprintf("%v-snapshot", id)
	res, err := r.Client.LRange(ctx, id, 0, 1).Result()
	if err != nil {
		return logger.AccountSnapshot{}, err
	}
	if len(res) == 0 {
		return logger.AccountSnapshot{}, nil
	}
	var result logger.AccountSnapshot
	if err := json.Unmarshal([]byte(res[0]), &result); err != nil {
		return logger.AccountSnapshot{}, err
	}
	return result, nil
}
