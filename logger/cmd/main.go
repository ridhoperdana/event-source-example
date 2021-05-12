package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/logger"

	"github.com/ridhoperdana/event-source-example/logger/internal"
)

func main() {
	redisClient := internal.NewRedisStorage(0, "", "localhost:6379")
	balanceWatermill := internal.NewBalanceWatermill("account-balance", "event-source-example", "pubsub-cred.json")
	loggerService := logger.NewService(redisClient)
	watermillHandler := internal.NewWatermillHandler("event-source-example", "gateway-account",
		"gateway-account-sub", "pubsub-cred.json", balanceWatermill, loggerService)
	watermillHandler.Process(context.TODO())
}
