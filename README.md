# Video Samples - extract resources from video

**🎥 Have you ever wanted to extract thumbnails or short videos from a real video? This project may be the answer you're looking for.**

See [the design docs](./DESIGN.md) if you want to know how it works.

## 🚀 Install

Coming soon

## 💡 Usage

Coming soon

## Development

```
docker compose up
```

This command will get up and running the following components:

* **Enqueuer** - enqueue jobs to extract the extractors
* **Worker** - workers to extract resources from video (using _ffmpeg_)
* **Redis** - used by _enqueuer_ and as a datastore to the extracted resources
* **Video API** - a dummy API that servers a JSON endpoint with video streamings URLs
* **Stream** - a live streaming generated with _ffmpeg_

## Tests

First, make sure you have the [ginkgo](http://onsi.github.io/ginkgo/) test runner installed.

Then run the test suite:

```
ginkgo -p -v ./...
```
