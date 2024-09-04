FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o gatewayApp /app/cmd/server/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gatewayApp .
COPY --from=builder /app/.envs/.local/config.env .

EXPOSE 5000
CMD ["./gatewayApp"]