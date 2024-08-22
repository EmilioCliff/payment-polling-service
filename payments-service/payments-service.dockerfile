FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o paymentApp /app/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/paymentApp .
COPY --from=builder /app/config.env .
COPY --from=builder /app/db/migrations /app/db/migrations

CMD ["./paymentApp"]