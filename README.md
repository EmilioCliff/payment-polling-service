# Go Microservices Project 🚀

Welcome to the **Go Microservices Project**! This repository contains a collection of microservices built with Go for practice. The project demonstrates a microservices architecture built with Go, Docker, and RabbitMQ. Each service is self-contained and serves a specific purpose in the larger architecture. Navigate to each service to learn more!

![Microservices Image](Microservice.drawio.png)

## Services Overview 🧩

- [Broker Service](./broker-service/README.md) - This is the gateway/entry point to our system. Clients interact with this service and the request is forwarded to the respective service through HTTP request or through a rabbit MQ queue.
- [Authentication Service](./authentication-service/README.md) - The authentication service connects to a postgres database and just authenticate the credentials received with the credentials saved in the database.
- [Logger Service](./logger-service/README.md) - The logger service is used to log important request data to a postgres database that is connected to.
- [Listener Service](./listener-service/README.md) - The listener is a dedicated service for setting up the rabbit queues and binds the routing key to the queues. It also redirects the message to respective services.
- [Payment Service](./payment-service/README.md) - The payment service is a experimental service used to learn how to intergrate the stripe card payment into the system.

## Tech Stack 🛠️

- **Go**: Backend services.
- **gRPC**: Service communication.
- **Docker**: Containerization.
- **RabbitMQ**: Message broker.

Explore the services by visiting their directories for more details.

## Getting Started 🛠️

To get started with the entire project, follow these steps:

1. **Clone the repository**

   ```bash
   git clone https://github.com/EmilioCliff/microservices-go
   ```

This next step require you to have docker compose installed and to get API keys from your [Stripe Account](.https://stripe.com/)

2. **Create .env files**

   At the root directory of the payment service create an .env file and fill the values below.

   ```
   STRIPE_SECRET_KEY=
   STRIPE_PUBLISHABLE_KEY=
   WEBHOOK=
   PORT=5000
   ```

You'll need to create a random user in the authentication postgres db.

3. **Run the services**

   ```bash
   docker compose up
   ```

4. **Making Request**
   We now cd into the frontend and start the sever with the following command

   ```
      cd ./front-end/front-end/cmd/web
      go run main.go
   ```

   We then go to `localhost:8082`

   With everything set up we can now start testing our services...😄

## Remarks 🤝

Happy Coding

## License 📝

This project is licensed under the MIT License.
