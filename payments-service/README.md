# Payments Service 🚀

Welcome to the **Payments Service**! This microservice is part of the larger Payment Polling System. It is designed to initate payments(withdrawal,deposits) to your Payd Wallet and retrieve/poll transaction details when requested.

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce http request
- **amqp091-go**: Used to connect and interact with rabbitmq for communications.
- **sqlc**: Used to generate type safe code from sql queries.
- **golang-migrate**: Used to apply updates to our database schema while helping in version control of our db schema.
- **gomock**: Used in mocking the sqlc generated postgres code, enables easy testing. [gomock](.https://github.com/uber-go/mock)
- **asynq**: Used in implementing scheduling tasks and retring failed process. [asynq](.https://github.com/hibiken/asynq)

## Configuration ⚙️

- Config setting: `PAYD_CALLBACK_URL` You will need to setup a callback url in the config file `./payments-service/.envs/.local/config.env`. The callback is used with payd to update transaction details after a successful transaction.

## Additional

Just as a side note. For temporary callback url you can check on [ngrok](.https://ngrok.com/). Please be sure it listens to the port your payment service is running on `:3030`  
The command would look like

```
    ngrok http 3030
```

## Changes

We could have configured the callback in the gateway service and prevent exposing our payment service to outside network but aah!! 🙃 we will work on that.
