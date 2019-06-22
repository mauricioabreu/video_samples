# thumbinator

<p align="center">
  <img src="https://github.com/mauricioabreu/thumbinator/raw/master/docs/thumbinator.png?raw=true" width="500">
  <p align="center">Generate thumbs from live streamings and videos on demand</p>
  <p align="center">
    <img src="https://travis-ci.org/mauricioabreu/thumbinator.svg?branch=master">
  </p>
</p>

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
* `make test` - Run test suite