package thumbinator

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func TestCollectThumbs(t *testing.T) {
	streams := []Stream{
		Stream{
			Name: "big_buck_bunny",
			URL:  "http://127.0.0.1:8080/play/hls/bunny/index.m3u8",
			TTL:  5,
		},
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	collector := Collector{Watcher: watcher, Store: dummyStore{}, Path: "../../../test/testdata/thumbnails"}
	go collector.CollectThumbs(streams)
	// Create a file system event so the collector wakes up
	thumbFile := "../../../test/testdata/thumbnails/000000001.jpg"
	err = ioutil.WriteFile(thumbFile, []byte("foobar"), 0644)
	if err != nil {
		panic(err)
	}
	// Give system time to sync write changes
	time.Sleep(50 * time.Millisecond)
	collector.Watcher.Close()
	err = os.Remove(thumbFile)
	if err != nil {
		panic(err)
	}

}
