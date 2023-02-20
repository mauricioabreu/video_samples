package extractor

import (
	"fmt"
)

type ThumbOptions struct {
	Input   string
	Output  string
	Scale   string
	Quality uint
}

func ExtractThumbs(title string, opts ThumbOptions, runner func(Command) error) error {
	args := []string{
		"-live_start_index", "-1",
		"-f", "hls",
		"-i", opts.Input,
		"-vf", "fps=1,scale=-1:360",
		"-vsync", "vfr",
		"-q:v", "5",
		"-threads", "1",
		fmt.Sprintf("%s/%s/%%09d.jpg", opts.Output, title),
	}
	return runner(Command{executable: "ffmpeg", args: args})
}
