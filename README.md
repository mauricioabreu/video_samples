# thumbinator

Generate thumbs from live streamings and videos on demand

![Architecture](/docs/thumbinator.png)

## Project goals

* Generate thumbs for live streaming videos
* Easy to deploy and run

## Commands
> Available make commands

* `make build-server` - Build base image to run the live streaming server
* `make server` - Run nginx server with nginx-ts module (responsible to produce HLS and DASH chunks)
* `make ingest` - Ingest a video to be handle by nginx-ts (big buck bunny strikes again)
* `make redis` - Run a redis instance to save all the generated thumbs
* `make run` - Run the thumbinator program