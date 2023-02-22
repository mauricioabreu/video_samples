create-video:
	ffmpeg -re -f lavfi -i "testsrc=size=1280x720:rate=30" \
		-pix_fmt yuv420p \
		-c:v libx264 -x264opts keyint=30:min-keyint=30:scenecut=-1 \
		-tune zerolatency -profile:v high -preset veryfast -bf 0 -refs 3 \
		-b:v 1400k -bufsize 1400k \
		-vf "drawtext=text='%{localtime}:box=1:fontcolor=black:boxcolor=white:fontsize=100':x=40:y=400'"  -hls_time 5 -hls_list_size 10 -hls_flags delete_segments -f hls testvideo/output.m3u8

clean-video:
	rm -f testvideo/*.m3u8
	rm -f testvideo/*.ts
