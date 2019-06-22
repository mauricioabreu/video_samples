package thumbinator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

const (
	sampleRate = 2
)

// Stream represent a stream to be processed
type Stream struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	TTL  int    `json:"ttl"`
}

type source interface {
	get() ([]Stream, error)
}

// JSONSource streams loaded from JSON files
type JSONSource struct {
	File string
}

// HTTPSource streams loaded from HTTP responses
type HTTPSource struct {
	URL     string
	Timeout time.Duration
}

func (js JSONSource) get() ([]Stream, error) {
	streams := make([]Stream, 0)
	data, err := ioutil.ReadFile(js.File)
	if err != nil {
		return streams, err
	}
	if err := json.Unmarshal(data, &streams); err != nil {
		return streams, err
	}

	return streams, nil
}

func (h HTTPSource) get() ([]Stream, error) {
	streams := make([]Stream, 0)
	netClient := http.Client{Timeout: h.Timeout}
	response, err := netClient.Get(h.URL)
	if err != nil {
		return streams, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return streams, err
	}
	if err := json.Unmarshal(data, &streams); err != nil {
		return streams, err
	}

	return streams, nil
}

// GetStreams retrieve a list of streams to process
func GetStreams(s source) ([]Stream, error) {
	return s.get()
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

// GenerateThumb start ffmpeg process to create thumbs
func GenerateThumb(streamingURL string, streamName string, path string) {
	createDir(filepath.Join(path, streamName))
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", fmt.Sprintf("%s", streamingURL), "-vf", "fps=1,scale=-1:360", "-vsync", "vfr", "-q:v", "5", "-threads", "1", fmt.Sprintf("%s/%s/%%09d.jpg", path, streamName)}
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

// CollectThumbs save all thumbs into redis
func CollectThumbs(streams []Stream, path string) {
	store := newRedisStore()
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
				seq, err := getSeqNumber(event.Name)
				if err != nil {
					log.Error(err)
					continue
				}
				if seq%sampleRate != 0 {
					if err := os.Remove(event.Name); err != nil {
						log.Debugf("Could not remove thumb file %s: %s", event.Name, err)
					}
					continue
				}

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
				stream, err := getStream(getStreamName(event.Name), streams)
				if err != nil {
					log.Error(err)
				}
				if err := store.SaveThumb(stream, timestamp, data); err != nil {
					log.Debugf("Could not save thumbs for %s: %s", stream.Name, err)
					continue
				}
				log.Debugf("Saved thumb for %s", stream.Name)
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

func getSeqNumber(filename string) (int, error) {
	onlyNumbersRegex := regexp.MustCompile("(\\d+)")
	seq, err := strconv.Atoi(onlyNumbersRegex.FindAllString(filepath.Base(filename), -1)[0])
	if err != nil {
		return 0, fmt.Errorf("could not convert %s to int", filename)
	}
	return seq, nil
}

func getStreamName(filename string) string {
	return filepath.Base(filepath.Dir(filename))
}

func getStream(streamName string, streams []Stream) (Stream, error) {
	for _, s := range streams {
		if s.Name == streamName {
			return s, nil
		}
	}
	return Stream{}, errors.New("stream not found")
}
