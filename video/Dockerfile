FROM alfg/ffmpeg
RUN apk add --no-cache python3
WORKDIR /video
COPY videoplay.sh /video
CMD [ "./videoplay.sh" ]
