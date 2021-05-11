# event-source-example
Trying out event sourcing

## Requirements
1. Google Pubsub
2. Redis

## How to Run
1. Setup 2 google pubsub topics:
  - account-balance
  - gateway-account

2. Setup 2 google pubsub subscriptions:
  - account-balance-sub
  - gateway-account-sub

3. Create Google Credential JSON key, and place it on root folder with `pubsub-cred.json` filename.
4. The app is consist of 3 binary:
 - gateway. This binary is to get data from client side, and send it to logger.
 - logger. This binary is to process any occured event (from gateway and balance)
 - balance. This binary is to process `Inputted Money` event and then process the data (change it into `Stored Money`) before send it back to logger.
5. Run the app with `go run logger/cmd/main.go` `go run balance/cmd/main.go` `go run gateway/cmd/main.go`

### Gateway App
Gateway app support 3 flag:
1. `account_id` the ID of the requesting account
2. `event_type` event that the account want to process. Currently support `INPUTTED_MONEY` and `STORED_MONEY`
3. `amount` amount of the processed money. Must be positive integer.
