package video_samples

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestCollectThumbs(t *testing.T) {
	streams := []Stream{
		{
			Name: "colors",
			URL:  "http://127.0.0.1:8080/play/hls/bunny/index.m3u8",
			TTL:  5,
		},
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	thumbsPath, _ := filepath.Abs("../../../test/testdata/thumbnails")
	ds := dummyStore{data: make(map[string][]byte, 0)}
	collector := Collector{
		Watcher: watcher,
		Store:   ds,
		Path:    thumbsPath,
	}
	go func() {
		collector.CollectThumbs(streams)
	}()
	time.Sleep(50 * time.Millisecond)
	// Create a file system event so the collector wakes up
	thumbFile, _ := filepath.Abs("../../../test/testdata/thumbnails/colors/000000002.jpg")
	file, err := os.OpenFile(thumbFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString("foobar")
	file.Sync()
	file.Close()
	time.Sleep(50 * time.Millisecond)

	collector.Watcher.Close()
	os.Remove(thumbFile)
	if string(ds.data["colors"]) != "foobar" {
		t.Errorf("Wrong thumb blob saved. Got %v wanted %v", string(ds.data["colors"]), "foobar")
	}
}
