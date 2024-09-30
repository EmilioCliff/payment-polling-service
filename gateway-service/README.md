# Gateway Service 🚀

Welcome to the **Gateway Service**! This microservice is part of the larger Payment Polling System. It is designed as an interface for clients and the services. Interacting with the client and then forwarding the request to the respective service.

## Endpoints ✨

`POST    /register` used to register a new user. Returns user created.
`POST     /login` used to login a user to a system. It returns the access_token used for protected endpoints.  
 `POST     /payments/initiate` used to initiate payments, can be withdrawal for withdrawing form your wallet or payments for depositing into your wallet. It return transaction_id which is used for checking on trabsaction status. 'PROTECTED=JWT'
`GET     /payments/status/:id` used to for polling transaction status. Returns transaction details. 'PROTECTED=JWT'

## Technologies Used 🛠️

- **gin-gonic**: A server to listen and produce json response.
- **amqp091-go**: Used to connect and interact with rabbitmq for communications.
- **gRPC**: Used to communicate with other services.
- **testcontainers**: Used to test rabbithandlers. Simplifies the mocking of receiving responses back and sending them without depending on the respective services to be running.
- **swaggo**: Used generate swagger2 specs. We later used [swagger2openapi](.https://www.npmjs.com/package/swagger2openapi) to change to openapi3, then updated some changes.
- **statik**: Used to serve static content which speeds up file serving since static files directly embedded into the binary, eliminating disk I/O straight [swagger2openapi](.https://github.com/rakyll/statik).

## Additional

Implemented different communication channels with the rest of the system, ie:

- authenication-service: Can communicate via gRPC call, http requests or rabbitmq queuing.  
  we can configure this in the `./gateway-service/internal/http/handlers.go` file

```
    // we can change the communication channel here ie
    // statusCode, rsp := s.RabbitService.RegisterUserViaRabbit(req)
    // statusCode, rsp := s.HTTPService.RegisterUserViaHttp(req)
    statusCode, rsp := s.GRPCService.RegisterUserViagRPC(req)
    if statusCode != http.StatusOK {
        ctx.JSON(statusCode, pkg.ErrorResponse(rsp.Message, rsp.StatusCode))

        return
    }
```

- payments-service: Can communicate only using rabbitmq.

**N/B:** After choosing the communication channel we will use in the handlers file, change the `./gateway-service/internal/http/handlers_test.go` to the respective channel so as to pass the test.

```
    // change here to the communication channel you are using
    s.GrpcService.RegisterUserViagRPCFunc = mockRegisterUserViagRPC
```
