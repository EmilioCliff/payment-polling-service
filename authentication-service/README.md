# Authentication Service 🚀

Welcome to the **Authentication Service**! This microservice is part of the larger Payment Polling System. It is designed to register new users and authenticate credentials provided by a client to the credentials stored in database.

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce http request
- **amqp091-go**: Used to connect and interact with rabbitmq for communications.
- **sqlc**: Used to generate type safe code from sql queries.
- **golang-migrate**: Used to apply updates to our database schema while helping in version control of our db schema.
- **gomock**: Used in mocking the sqlc generated postgres code, enables easy testing. [gomock](.https://github.com/uber-go/mock)

## Additional

Just as a side note. It opens a grpc server that is used to give user information/data by the payment service during the initiate payment process.
