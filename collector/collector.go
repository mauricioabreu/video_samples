package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mauricioabreu/video_samples/collector/watcher"
	"github.com/rs/zerolog/log"
)

type File struct {
	Path    string
	Dir     string
	Data    []byte
	ModTime int64
}

func Collect(path string) {
	files, err := watcher.Watch(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize collector")
	}

	for file := range files {
		log.Debug().Msgf("File found: %s", file)
	}
}

func NewFile(path string) (*File, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file info: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}
	return &File{
		Path:    path,
		Dir:     filepath.Dir(path),
		Data:    data,
		ModTime: fileInfo.ModTime().Unix(),
	}, nil
}
