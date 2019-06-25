build-server:
	docker build -t nginx-live .

server: build-server
	docker run -it -p 1935:1935 -p 8080:8080 nginx-live

ingest:
	docker run --net="host" --rm -v $(shell pwd):/files jrottenberg/ffmpeg:4.1 -re -i /files/big_buck_bunny_480p.mp4 -c copy -f mpegts http://127.0.0.1:8080/publish/bunny

redis:
	docker run -d -p 6379:6379 redis

run:
	go run main.go generate

test:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

edge:
	go run main.go server

install:
	go install

.PHONY: build-server server ingest redis run test edge