package main

import (
	"context"

	"github.com/ridhoperdana/event-source-example/logger"

	"github.com/ridhoperdana/event-source-example/logger/internal"
)

func main() {
	redisClient := internal.NewRedisStorage(0, "", "localhost:6379")
	balanceWatermill := internal.NewBalanceWatermill("account-balance", "event-source-example",
		"pubsub-cred.json", "gateway-account")
	marketPlaceWatermill := internal.NewMarketplaceWatermill("account-marketplace", "event-source-example",
		"pubsub-cred.json", "gateway-marketplace")
	loggerService := logger.NewService(redisClient)
	watermillHandler := internal.NewWatermillHandler(balanceWatermill, marketPlaceWatermill, loggerService)
	watermillHandler.Process(context.TODO())
}
