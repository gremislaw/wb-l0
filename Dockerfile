#Docker Pipeline
FROM golang:1.23.2-alpine3.20 as builder
WORKDIR /app
COPY . /app
RUN go mod download && \
    go build -o ./bin/order_service ./cmd/app/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/order_service ./bin/order_service
EXPOSE 9000
ENTRYPOINT ["./bin/order_service"]
