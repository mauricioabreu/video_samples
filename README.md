# thumbinator

<p align="center">
  <img src="https://github.com/mauricioabreu/thumbinator/raw/master/docs/thumbinator.png?raw=true" width="500">
  <p align="center">Generate thumbs from live streamings and videos on demand</p>
  <p align="center">
    <a href="https://travis-ci.org/mauricioabreu/thumbinator">
      <img src="https://travis-ci.org/mauricioabreu/thumbinator.svg?branch=master">
    </a>
    <a href="https://codecov.io/gh/mauricioabreu/thumbinator">
      <img src="https://codecov.io/gh/mauricioabreu/thumbinator/branch/master/graph/badge.svg">
    </a>
  </p>
</p>

## Project goals

* Generate thumbs for live streaming videos
* Easy to deploy and run

## Installing

Currently we only support installing `thumbinator` from source:

```console
git clone git@github.com:mauricioabreu/thumbinator.git
cd thumbinator
make install
```

## Trying it out!

This project comes with tools to try it locally, without having a real live streaming on the internet.
To achieve it we use some open source technologies like `ffmpeg` and `nginx-ts`.

First, we need to start a `nginx` server to produce our streaming playlists:
```
make server
```

**p.s**: first run will take some time because it will compile and generate an image with `nginx-ts`.

Now that we have a streaming server we can `ingest` a video. You can use the well known [Big Buck Bunny](https://peach.blender.org/download/)

```
make ingest
```

[VLC](https://www.videolan.org/vlc/) can be used to test if our streaming works.

## Extracting thumbs

Regardless of having local or online live streamings, we use `thumbinator` to generate thumbs.
We have a `streams.json` file that serves as example to determine which URLs will be consumed to produce our thumbs.

`thumbinator` comes with two entrypoints: `generate` and `server`

### Generate

Start a program to extract, collect and save thumbs on `redis`:

```console
thumbinator generate
```

### Server

Start an HTTP server to query the saved thumbs for the given stream:

```console
thumbinator server
```

For both commands you can get help:

```console
thumbinator help <command>
```

## Commands
> Available make commands

* `make build-server` - Build base image to run the live streaming server
* `make server` - Run nginx server with nginx-ts module (responsible to produce HLS and DASH chunks)
* `make ingest` - Ingest a video to be handle by nginx-ts (big buck bunny strikes again)
* `make redis` - Run a redis instance to save all the generated thumbs
* `make run` - Run the thumbinator program
* `make test` - Run test suite
* `make edge` - Run an HTTP server to delivery our nice thumbnails