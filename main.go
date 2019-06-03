package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	streamingURL := "http://127.0.0.1:8080/play/hls/bunny/index.m3u8"
	generateThumb(streamingURL)
	fmt.Println("Generating thumbnail...")
}

func generateThumb(streamingURL string) {
	args := []string{"-live_start_index", "-1", "-f", "hls", "-i", fmt.Sprintf("%s", streamingURL), "-vf", "fps=1,scale=-1:169", "-vsync", "vfr", "-q:v", "5", "-threads", "1", "%09d.jpg"}
	cmd := exec.Command("ffmpeg", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ffmpeg outout: %q\n", out.String())
}
