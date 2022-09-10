#!/bin/sh

set -x

python3 -m http.server 8080 --bind 0.0.0.0 &

ffmpeg -hide_banner \
    -re -f lavfi -i "smptehdbars=size=1280x720:rate=30,format=yuv420p" \
    -vf "drawtext=fontsize=96:fontcolor=white:text='%{localtime\:%T}':fontfile='OpenSans-Bold.ttf'" \
    -map 0:v:0 -map 0:v:0 -map 0:v:0 \
    -c:v libx264 -preset ultrafast -tune zerolatency -profile:v high \
    -b:v:0 300k -s:v:0 480:360 -bufsize 700k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -b:v:1 700k -s:v:1 640:480 -bufsize 1500k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -b:v:2 1400k -s:v:2 1280:760 -bufsize 2800k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -var_stream_map "v:0,name:360p v:1,name:480p v:2,name:720p" \
    -hls_list_size 10 -threads 0 -f hls \
    -hls_init_time 5 -hls_time 5 \
    -master_pl_name "colors.m3u8" -y "colors-%v.m3u8"
