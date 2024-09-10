FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o paymentApp /app/cmd/server/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/paymentApp .
COPY --from=builder /app/.envs/.local/config.env .
COPY --from=builder /app/internal/postgres/migrations  /app/migrations

EXPOSE 3030

CMD ["./paymentApp"]