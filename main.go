package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func main() {
	thumbsPath := "thumbnails"
	streamingURL := "http://127.0.0.1:8080/play/hls/bunny/index.m3u8"
	generateThumb(streamingURL)
	collectThumbs(thumbsPath)
}

func generateThumb(streamingURL string) {
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", fmt.Sprintf("%s", streamingURL), "-vf", "fps=1,scale=-1:169", "-vsync", "vfr", "-q:v", "5", "-threads", "1", "thumbnails/%09d.jpg"}
	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	log.Printf("Generating thumbnail for %s\n", streamingURL)
}

func collectThumbs(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for err := range watcher.Errors {
			log.Println(err)
		}
	}()

	if err := watcher.Add(path); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		for event := range watcher.Events {
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Printf("Received file %s", event.Name)
			}
		}
		done <- true
	}()

	select {
	case <-done:
		log.Print("Done watching files...")
	}
}
