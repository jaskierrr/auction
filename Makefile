.PHONY: run

run:
		go run cmd/auction/main.go

build:
		docker-compose up --build -d

logs:
		docker logs -f auction_app

test:
		go test -v -count=1 ./...
