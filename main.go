package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

type stream struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type source interface {
	get() ([]stream, error)
}

// JSONSource sources loaded from JSON files
type JSONSource struct {
	file string
}

func (js JSONSource) get() ([]stream, error) {
	streams := make([]stream, 0)
	data, err := ioutil.ReadFile(js.file)
	if err != nil {
		return streams, err
	}
	if err := json.Unmarshal(data, &streams); err != nil {
		return streams, err
	}

	return streams, nil
}

func getStreams(s source) ([]stream, error) {
	return s.get()
}

func main() {
	thumbsPath := "thumbnails"
	streams, err := getStreams(JSONSource{file: "streams.json"})
	if err != nil {
		log.Fatalf("Could not retrieve streams to process: %s", err)
	}
	log.Infof("Retrieved the following streams: %+v", streams)
	for _, s := range streams {
		log.Debugf("Generating thumbs for %s with URL %s", s.Name, s.URL)
		generateThumb(s.URL, s.Name, thumbsPath)
	}
	collectThumbs(thumbsPath)
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func generateThumb(streamingURL string, streamName string, path string) {
	createDir(filepath.Join(path, streamName))
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", fmt.Sprintf("%s", streamingURL), "-vf", "fps=1,scale=-1:169", "-vsync", "vfr", "-q:v", "5", "-threads", "1", fmt.Sprintf("%s/%s/%%09d.jpg", path, streamName)}
	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func getSubDirs(sourcePath string) []string {
	var paths []string
	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	return paths
}

func collectThumbs(path string) {
	client := newClient()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watcher.Close()

	go func() {
		for err := range watcher.Errors {
			log.Error(err)
		}
	}()

	if err := watcher.Add(path); err != nil {
		log.Fatal(err)
	}
	for _, d := range getSubDirs(path) {
		if err := watcher.Add(d); err != nil {
			log.Fatal(err)
		}
	}

	done := make(chan bool)
	go func() {
		for event := range watcher.Events {
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Debugf("Received event: %s", event.Name)
				pathInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Fatalf("Could not read path metadata for %s: %s", event.Name, err)
				}
				if pathInfo.IsDir() {
					watcher.Add(event.Name)
					continue
				}

				data, err := ioutil.ReadFile(event.Name)
				if err != nil {
					log.Errorf("Could not read file: %s", err)
					continue
				}

				timestamp := pathInfo.ModTime().UTC().Unix()
				streamName := getStreamName(event.Name)
				if err := saveThumb(client, streamName, timestamp, data); err != nil {
					log.Debugf("Could not save thumbs for %s: %s", streamName, err)
					continue
				}
				log.Debugf("Saved thumb for %s", streamName)
				if err := os.Remove(event.Name); err != nil {
					log.Debugf("Could not remove thumb file %s: %s", event.Name, err)
				}
			}
		}
		done <- true
	}()

	select {
	case <-done:
		log.Info("Done watching files...")
	}
}

func getStreamName(filename string) string {
	return filepath.Base(filepath.Dir(filename))
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

	if err = c.Set(fmt.Sprintf("thumbs/blob/%s", id.String()), blob, time.Duration(60)*time.Second).Err(); err != nil {
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
	log.Info(pong, err)
	return client
}
