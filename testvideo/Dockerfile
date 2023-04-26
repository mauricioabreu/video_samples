FROM python:3.9-alpine

WORKDIR /testvideo

RUN apk upgrade -U && apk add ffmpeg

COPY stream.sh OpenSans-Bold.ttf /testvideo/

RUN chmod +x /testvideo/stream.sh
