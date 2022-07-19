run:
	docker compose up

test:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	go build -o video_samples

install:
	go install

.PHONY: run test build install
