FROM golang:1.18-alpine

RUN apk add --no-cache gcc \
    libc-dev \
    linux-headers \
    ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /video_samples
