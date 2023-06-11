# Run project
run:
	docker compose up

# Run tests
test:
	ginkgo -p -v ./...

# Purge generated video assets
clean-video:
	rm -f testvideo/*.m3u8
	rm -f testvideo/*.ts

# Purge generated thumbs
clean-thumbs:
	rm -rf testvideo/thumbs/**/*.jpg
