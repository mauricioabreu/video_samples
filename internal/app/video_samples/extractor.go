package video_samples

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

const (
	sampleRate = 2
)

type Collector struct {
	Watcher   *fsnotify.Watcher
	Store     store
	Path      string
	Recursive bool
}

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
	data, err := os.ReadFile(js.File)
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
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
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
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Error().Err(err).Msg("")
		}
	}
}

// GenerateThumb start ffmpeg process to create thumbs
func GenerateThumb(streamingURL string, streamName string, path string) {
	createDir(filepath.Join(path, streamName))
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", streamingURL, "-vf", "fps=1,scale=-1:360", "-vsync", "vfr", "-q:v", "5", "-threads", "1", fmt.Sprintf("%s/%s/%%09d.jpg", path, streamName)}
	cmd := exec.Command("ffmpeg", args...)
	log.Debug().Msgf("Executing ffmpeg with args: %v", args)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func getSubDirs(sourcePath string) []string {
	var paths []string
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return paths
}

// CollectThumbs save all thumbs into redis
func (c Collector) CollectThumbs(streams []Stream) {
	go func() {
		for err := range c.Watcher.Errors {
			log.Error().Err(err).Msg("")
		}
	}()

	if err := c.Watcher.Add(c.Path); err != nil {
		log.Fatal().Err(err).Msg("")
	}
	for _, d := range getSubDirs(c.Path) {
		if err := c.Watcher.Add(d); err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}

	done := make(chan bool)
	go func() {
		for event := range c.Watcher.Events {
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Debug().Msgf("Received event: %s", event.Name)
				seq, err := getSeqNumber(event.Name)
				if err != nil {
					log.Error().Err(err).Msg("")
					continue
				}
				if seq%sampleRate != 0 {
					if err := os.Remove(event.Name); err != nil {
						log.Debug().Msgf("Could not remove thumb file %s", event.Name)
					}
					continue
				}

				pathInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Fatal().Err(err).Msgf("Could not read path metadata for %s", event.Name)
				}
				if pathInfo.IsDir() {
					if err := c.Watcher.Add(event.Name); err != nil {
						log.Error().Err(err).Msg("")
					}
					continue
				}

				data, err := os.ReadFile(event.Name)
				if err != nil {
					log.Error().Err(err).Msg("Could not read file")
					continue
				}

				if len(data) == 0 {
					log.Info().Msgf("File %s is empty", event.Name)
					continue
				}

				timestamp := pathInfo.ModTime().UTC().Unix()
				stream, err := getStream(getStreamName(event.Name), streams)
				if err != nil {
					log.Error().Err(err).Msg("")
				}
				if err := c.Store.SaveThumb(stream, timestamp, data); err != nil {
					log.Debug().Msgf("Could not save thumbs for %s", stream.Name)
					continue
				}
				log.Debug().Msgf("Saved thumb for %s", stream.Name)
				if err := os.Remove(event.Name); err != nil {
					log.Debug().Msgf("Could not remove thumb file %s", event.Name)
				}
			}
		}
		done <- true
	}()

	<-done
	log.Info().Msg("Done watching files...")
}

func getSeqNumber(filename string) (int, error) {
	onlyNumbersRegex := regexp.MustCompile(`(\d+)`)
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
