# Broker Service 🚀

Welcome to the **Broker Service**! This microservice is part of the larger Go Microservices project. It is designed as an interface for our clients and the services. Interacting with the client and then forwarding the request to the respective service.

## Endpoints ✨

- Feature 1
  `GET    /ping` used to test if the service is up and healthy
  `POST     /handler` used to receive the clients request and route to respective services using either the rabbit queue or http request
  `POST     /webhook` used together with the payment api for stripe callbacks

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce http request.
- **amqp091-go**: Used to for interaction with the AMQP that being RabbitMQ where request are queued.

## Configuration ⚙️

- Config setting 1: `DB_URL` Set in the compose file
- Confiq setting 2: `STRIPE_SECRET_KEY=` Provided by Stripe
- Confiq setting 3: `STRIPE_PUBLISHABLE_KEY=` Provided by Stripe
- Confiq setting 4: `WEBHOOK=` Used a live callback url. I used [ngrok](.https://ngrok.com/)
- Confiq setting 5: `PORT=5000`
