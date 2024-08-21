FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o authApp /app/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/authApp .
COPY --from=builder /app/config.env .
COPY --from=builder /app/api/server.go /app/api/server.go
COPY --from=builder /app/db/migrations /app/db/migrations
COPY --from=builder /app/utils/* /app/utils/

EXPOSE 5000
EXPOSE 5050

CMD ["./authApp"]