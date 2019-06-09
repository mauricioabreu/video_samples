package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
)

func main() {
	thumbsPath := "thumbnails"
	streamingURL := "http://127.0.0.1:8080/play/hls/bunny/index.m3u8"
	generateThumb(streamingURL, thumbsPath)
	collectThumbs(thumbsPath)
}

func generateThumb(streamingURL string, path string) {
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", fmt.Sprintf("%s", streamingURL), "-vf", "fps=1,scale=-1:169", "-vsync", "vfr", "-q:v", "5", "-threads", "1", fmt.Sprintf("%s/%%09d.jpg", path)}
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
				if err := os.Remove(event.Name); err != nil {
					log.Printf("Could not remove thumb file %s: %s", event.Name, err)
				}
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
		return err
	}

	err = c.ZAdd(fmt.Sprintf("thumbs/%s", key), redis.Z{Score: float64(timestamp), Member: id.String()}).Err()
	if err != nil {
		return err
	}

	if err = c.Set(id.String(), blob, time.Duration(60)*time.Second).Err(); err != nil {
		return err
	}

	if err := c.ZRemRangeByScore(fmt.Sprintf("thumbs/%s", key), "-inf", strconv.Itoa(int(timestamp)-60)).Err(); err != nil {
		return err
	}

	return nil
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
