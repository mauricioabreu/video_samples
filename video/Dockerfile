FROM linuxserver/ffmpeg:latest
RUN apt update && apt-get -y install python3 tini
WORKDIR /video
COPY videoplay.sh /video
CMD [ "tini", "--", "./videoplay.sh" ]
