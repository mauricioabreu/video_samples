FROM golang:1.20

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /video_samples
