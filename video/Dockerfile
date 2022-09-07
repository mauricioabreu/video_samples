FROM alfg/ffmpeg
RUN apk add --no-cache python3 tini
WORKDIR /video
COPY videoplay.sh /video
CMD [ "tini", "--", "./videoplay.sh" ]
