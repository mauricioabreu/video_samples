package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Path    string
	Dir     string
	Data    []byte
	ModTime int64
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
