build-server:
	docker build -t nginx-live .

server: build-server
	docker run -it -p 1935:1935 -p 8080:8080 nginx-live

ingest:
	docker run --net="host" --platform linux/amd64 --rm -v $(shell pwd):/files jrottenberg/ffmpeg -hide_banner \
    -re -f lavfi -i "testsrc2=size=1280x720:rate=30,format=yuv420p" \
    -f lavfi -i "sine=frequency=220:beep_factor=4:duration=5" \
    -c:v libx264 -preset ultrafast -tune zerolatency -profile:v high \
    -b:v 1400k -bufsize 2800k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -c:a aac -b:a 128k \
    -window_size 5 -extra_window_size 10 -remove_at_exit 1 -adaptation_sets "id=0,streams=v id=1,streams=a" \
    -fflags +genpts \
    -f mpegts http://127.0.0.1:8080/publish/colors

redis:
	docker run -d -p 6379:6379 redis

run:
	go run main.go generate

test:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

edge:
	go run main.go server

build:
	go build -o video_samples

install:
	go install

.PHONY: build-server server ingest redis run test edge
