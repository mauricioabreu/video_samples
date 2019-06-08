package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
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
	client := newClient()
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
				data, err := ioutil.ReadFile(event.Name)
				if err != nil {
					log.Printf("Could not read file: %s", err)
					continue
				}
				pathInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Fatalf("Could not read path metadata for %s: %s", event.Name, err)
				}
				timestamp := pathInfo.ModTime().UTC().Unix()
				if err := saveThumb(client, "big_buck_bunny", timestamp, data); err != nil {
					log.Printf("Could not save thumbs for %s: %s", "big_buck_bunny", err)
					continue
				}
				log.Printf("Saved thumb for %s", "big_buck_bunny")
			}
		}
		done <- true
	}()

	select {
	case <-done:
		log.Print("Done watching files...")
	}
}

func saveThumb(c *redis.Client, key string, timestamp int64, blob []byte) error {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	return c.ZAdd(fmt.Sprintf("thumbs/%s", key), redis.Z{Score: float64(timestamp), Member: id}).Err()
}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	return client
}
