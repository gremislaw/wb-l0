.PHONY: all test build format up down stop

all: build

run:
	cd ./ && go run cmd/app/main.go

build:
	go mod download
	go build -o ./bin/order_service

test:
	go test ./...

up:
	sudo docker compose up -d --build
	sudo docker logs -tf order_service

stop:
	sudo docker compose stop

down:
	sudo docker compose down

rebuild: clean
	go mod tidy
	make build

format:
	go fmt ./...

clean:
	rm -rf ./bin
	rm -rf ./data