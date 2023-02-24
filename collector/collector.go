package collector

import "github.com/mauricioabreu/video_samples/collector/watcher"

func Collect(path string) {
	watcher.Watch(path)
}
